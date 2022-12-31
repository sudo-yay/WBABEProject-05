package order

import (
	"fmt"
	"github.com/codestates/WBABEProject-05/common/util"
	"github.com/codestates/WBABEProject-05/model/entity"
	"github.com/codestates/WBABEProject-05/model/menu"
	"github.com/codestates/WBABEProject-05/model/receipt"
	error2 "github.com/codestates/WBABEProject-05/protocol/error"
	"github.com/codestates/WBABEProject-05/protocol/page"
	"github.com/codestates/WBABEProject-05/protocol/request"
	"github.com/codestates/WBABEProject-05/protocol/response"
	"time"
)

type orderRecordService struct {
	receiptModel receipt.ReceiptModeler
	menuModel    menu.MenuModeler
}

var instance *orderRecordService

func NewOrderRecordService(rd receipt.ReceiptModeler, md menu.MenuModeler) *orderRecordService {
	if instance != nil {
		return instance
	}

	instance = &orderRecordService{receiptModel: rd, menuModel: md}
	return instance
}

func (o *orderRecordService) RegisterOrderRecord(order *request.RequestOrder) (string, error) {
	rct, err := order.ToNewReceipt()
	if err != nil {
		return "", err
	}

	menus, err := o.menuModel.SelectMenusByIds(order.StoreId, order.Menus)
	if err != nil {
		return "", err
	}

	totalPrice := o.sumMenusPrice(menus)
	rct.Price = totalPrice

	toDayCnt, err := o.receiptModel.SelectToDayTotalCount()
	if err != nil {
		return "", err
	}
	numbering := fmt.Sprintf("%d-%d", time.Now().UnixNano(), toDayCnt)
	rct.Numbering = numbering

	insertedId, err := o.receiptModel.InsertReceipt(rct)
	if err != nil {
		return "", err
	}

	return insertedId, nil
}
func (o *orderRecordService) ModifyOrderRecordFromCustomer(order *request.RequestPutCustomerOrder) (string, error) {
	foundOrder, err := o.receiptModel.SelectReceiptByID(order.ID)
	if err != nil {
		return "", err
	}
	if foundOrder.Status != entity.Waiting {
		return "", error2.AlreadyReceivedOrderError.New()
	}
	if foundOrder.CustomerID.Hex() != order.CustomerID {
		return "", error2.BadAccessOrderError.New()
	}

	// TODO 취소시키고 새로접수하는게 맞는것 같다.
	if _, err := o.receiptModel.UpdateCancelReceipt(foundOrder); err != nil {
		return "", err
	}

	savedID, err := o.RegisterOrderRecord(order.ToRequestOrder())
	if err != nil {
		return "", err
	}

	return savedID, nil
}

func (o *orderRecordService) ModifyOrderRecordFromStore(order *request.RequestPutStoreOrder) (int, error) {
	foundOrder, err := o.receiptModel.SelectReceiptByID(order.ID)
	if err != nil {
		return 0, err
	}
	foundOrder.Status = order.Status
	updatedCnt, err := o.receiptModel.UpdateReceiptStatus(foundOrder)
	if err != nil {
		return 0, err
	}
	return updatedCnt, nil
}

func (o *orderRecordService) FindOrderRecordsSortedPage(userID string, pg *request.RequestPage) (*page.PageData[any], error) {
	skip := pg.CurrentPage * pg.ContentCount
	if skip > 0 {
		skip--
	}

	receipts, err := o.receiptModel.SelectSortLimitedReceipt(userID, pg.Sort, skip, pg.ContentCount)
	if err != nil {
		return nil, err
	}

	totalCount, err := o.receiptModel.SelectTotalCount(userID)
	if err != nil {
		return nil, err
	}

	pgInfo := pg.NewPageInfo(totalCount)

	return page.NewPageData(receipts, pgInfo), nil
}

func (o *orderRecordService) SelectReceipts() {

}

func (o *orderRecordService) FindOrderRecord(orderID string) (*response.ResponseOrder, error) {
	foundReceipt, err := o.receiptModel.SelectReceiptByID(orderID)
	if err != nil {
		return nil, err
	}
	menuIDs := util.ConvertObjIDsToStrings(foundReceipt.Menus)
	menus, err := o.menuModel.SelectMenusByIds(foundReceipt.StoreID.Hex(), menuIDs)
	if err != nil {
		return nil, err
	}
	resOrder := response.FromReceiptAndMenus(foundReceipt, menus)
	return resOrder, nil
}

func (o *orderRecordService) FiendSelectedMenusTotalPrice(storeID string, menuIDs []string) (*response.ResponseCheckPrice, error) {
	menus, err := o.menuModel.SelectMenusByIds(storeID, menuIDs)
	if err != nil {
		return nil, err
	}
	totalPrice := o.sumMenusPrice(menus)
	resCheckPrice := response.NewResponseCheckPrice(menus, totalPrice)
	return resCheckPrice, nil
}

func (o *orderRecordService) sumMenusPrice(menus []*entity.Menu) int {
	var totalPrice int
	for _, menu := range menus {
		totalPrice += menu.Price
	}
	return totalPrice
}
