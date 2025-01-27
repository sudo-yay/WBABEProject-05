package customer

import (
	"fmt"
	"github.com/codestates/WBABEProject-05/common/convertor"
	"github.com/codestates/WBABEProject-05/common/enum"
	error2 "github.com/codestates/WBABEProject-05/common/error"
	"github.com/codestates/WBABEProject-05/common/util"
	"github.com/codestates/WBABEProject-05/logger"
	"github.com/codestates/WBABEProject-05/model/entity"
	"github.com/codestates/WBABEProject-05/model/menu"
	"github.com/codestates/WBABEProject-05/model/query"
	"github.com/codestates/WBABEProject-05/model/receipt"
	"github.com/codestates/WBABEProject-05/model/review"
	"github.com/codestates/WBABEProject-05/protocol/page"
	"github.com/codestates/WBABEProject-05/protocol/request"
	"math"
)

var instance *menuReviewService

type menuReviewService struct {
	reviewModel  review.ReviewModeler
	menuModel    menu.MenuModeler
	receiptModel receipt.ReceiptModeler
}

func NewMenuReviewService(rMod review.ReviewModeler, mMod menu.MenuModeler, receiptMod receipt.ReceiptModeler) *menuReviewService {
	if instance != nil {
		return instance
	}
	instance = &menuReviewService{
		reviewModel:  rMod,
		menuModel:    mMod,
		receiptModel: receiptMod,
	}
	return instance
}

func (m *menuReviewService) FindReviewSortedPageByMenuID(menuID string, pg *request.RequestPage) (*page.PageData[any], error) {
	skip := util.NewSkipNumber(pg.CurrentPage, pg.ContentCount)

	pageQuery := query.NewPageQuery(pg.Sort.Name, pg.Sort.Direction, skip, pg.ContentCount)

	reviews, err := m.reviewModel.SelectSortLimitedReviewsByMenuID(menuID, pageQuery)
	if err != nil {
		return nil, err
	}

	totalCount, err := m.reviewModel.SelectTotalCountByMenuID(menuID)
	if err != nil {
		return nil, err
	}

	pgInfo := pg.ToPageInfo(int(totalCount))

	return page.NewPageData(reviews, pgInfo), nil
}

func (m *menuReviewService) FindReviewSortedPageByUserID(ID, userRole string, pg *request.RequestPage) (*page.PageData[any], error) {
	skip := util.NewSkipNumber(pg.CurrentPage, pg.ContentCount)

	pageQuery := query.NewPageQuery(pg.Sort.Name, pg.Sort.Direction, skip, pg.ContentCount)

	reviews, err := m.reviewModel.SelectSortLimitedReviewsByUserID(ID, userRole, pageQuery)
	if err != nil {
		return nil, err
	}

	totalCount, err := m.reviewModel.SelectTotalCountByUserID(ID, userRole)
	if err != nil {
		return nil, err
	}

	pgInfo := pg.ToPageInfo(int(totalCount))

	return page.NewPageData(reviews, pgInfo), nil
}

func (m *menuReviewService) RegisterMenuReview(review *request.RequestPostReview) (string, error) {
	foundReceipt, err := m.receiptModel.SelectReceiptByID(review.OrderID)
	if err != nil {
		return enum.BlankSTR, error2.DoesNotExistsOrderErr
	}

	if s, err := m.checkPossiblePostReview(review, foundReceipt); err != nil {
		return s, err
	}

	newR, err := review.ToPostReview()
	if err != nil {
		return enum.BlankSTR, err
	}

	savedID, err := m.reviewModel.InsertReview(newR)
	if err != nil {
		return enum.BlankSTR, err
	}

	// Rating 은 비즈니스상 중요하지않아보여 따로 컨틀롤하지 않는 고루틴 활용
	go m.updateMenuScores(review)

	return savedID, nil
}

func (m *menuReviewService) checkPossiblePostReview(review *request.RequestPostReview, foundReceipt *entity.Receipt) (string, error) {
	foundMIDs := convertor.ConvertOBJIDsToStrings(foundReceipt.MenuIDs)

	existMap := util.ConvertSliceToExistMap(foundMIDs)

	if _, exists := existMap[review.MenuID]; !exists {
		return enum.BlankSTR, error2.DoesNotPostReviewErr
	}

	if foundReceipt.CustomerID.Hex() != review.CustomerID || foundReceipt.StoreID.Hex() != review.StoreID {
		return enum.BlankSTR, error2.DoesNotPostReviewErr
	}

	return "", nil
}

func (m *menuReviewService) updateMenuScores(review *request.RequestPostReview) {
	foundM, err := m.menuModel.SelectMenuByID(review.MenuID)
	if err != nil {
		logger.AppLog.Error(err.Error())
	}

	foundM.TotalReviewScore += review.Rating
	foundM.ReviewCount++
	foundM.Rating = math.Round((float64(foundM.TotalReviewScore)/float64(foundM.ReviewCount))*10) / 10

	rating, err := m.menuModel.UpdateMenuRating(foundM)
	if err != nil || rating == 0 {
		MSG := fmt.Sprintf("does not update rating Menu ID %v", foundM.ID)
		logger.AppLog.Error(MSG)
	}
}
