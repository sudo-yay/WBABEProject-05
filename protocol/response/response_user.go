package response

import (
	"github.com/codestates/WBABEProject-05/model/entity"
)

type ResponseUser struct {
	ID               string           `json:"user_id"`
	Name             string           `json:"name"`
	NicName          string           `json:"nic_name"`
	PhoneNumber      string           `json:"phone_number"`
	Role             string           `json:"role"`
	OrderAddr        *ResponseAddress `json:"pre_order_addr"`
	OrderPhoneNumber string           `json:"pre_order_phone_number"`
}

func NewResponseUserFromUser(user *entity.User) *ResponseUser {
	return &ResponseUser{
		ID:               user.ID.Hex(),
		Name:             user.Name,
		NicName:          user.NicName,
		PhoneNumber:      user.PhoneNumber,
		Role:             user.Role,
		OrderAddr:        FromAddr(user.PreOrderInfo.Address),
		OrderPhoneNumber: user.PhoneNumber,
	}
}
