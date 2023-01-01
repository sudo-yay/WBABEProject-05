package entity

import (
	"github.com/codestates/WBABEProject-05/model/entity/dom"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StoreID    primitive.ObjectID `bson:"store_id,omitempty"`
	CustomerID primitive.ObjectID `bson:"customer_id,omitempty"`
	Menu       primitive.ObjectID `bson:"menu_id,omitempty"`
	Content    string             `bson:"content,omitempty"`
	Rating     int                `bson:"rating"`
	BaseTime   *dom.BaseTime      `bson:"base_time"`
}
