package request

import (
	"github.com/codestates/WBABEProject-05/model/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestPutUser struct {
	Name        string `json:"name" validate:"required"`
	NicName     string `json:"nic_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Role        string `json:"role"  validate:"required, eq=store|eq=user"`
}

func (r *RequestPutUser) NewUpdatePutUser(ID string) (*entity.User, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:          objID,
		Name:        r.Name,
		NicName:     r.NicName,
		PhoneNumber: r.PhoneNumber,
		Role:        r.Role,
	}, nil
}
