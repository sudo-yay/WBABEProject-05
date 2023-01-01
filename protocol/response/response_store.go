package response

import (
	"github.com/codestates/WBABEProject-05/model/entity"
	"github.com/codestates/WBABEProject-05/model/util"
)

type ResponseStore struct {
	StoreID        string           `json:"store_id"`
	UserID         string           `json:"user_id"`
	Name           string           `json:"name"`
	Address        *ResponseAddress `json:"address"`
	StorePhone     string           `json:"store_phone"`
	RecommendMenus []string         `json:"recommend_menus"`
}

func NewResponseStore(store *entity.Store) *ResponseStore {
	recommendIDs := util.ConvertOBJIDsToStrings(store.RecommendMenus)
	return &ResponseStore{
		StoreID:        store.ID.Hex(),
		UserID:         store.UserID.Hex(),
		Name:           store.Name,
		Address:        FromAddr(store.Address),
		StorePhone:     store.StorePhone,
		RecommendMenus: recommendIDs,
	}
}
