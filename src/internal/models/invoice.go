package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      string             `json:"userID" bson:"userID"`
	InvTitle    string             `json:"invTitle" bson:"invTitle"`
	InvNum      int                `json:"invNum" bson:"invNum"`
	CreatedTime string             `json:"createdTime" bson:"createdTime"`
	Balance     float64            `json:"balance" bson:"balance"`
	Notes       string             `json:"notes" bson:"notes"`
	Dispatch    bool               `json:"dispatch" bson:"dispatch"`
	Discount    bool               `json:"discount" bson:"discount"`
	ColorLine   string             `json:"colorLine" bson:"colorLine"`
	Currency    string             `json:"currency" bson:"currency"`
	From        From               `json:"from" bson:"from"`
	To          To                 `json:"to" bson:"to"`
	InvList     []InvItem          `json:"invList" bson:"invList"`
}

type From struct {
	Name          string  `json:"name" bson:"name"`
	EmailFrom     string  `json:"emailFrom" bson:"emailFrom"`
	Address       Address `json:"address" bson:"address"`
	Phone         string  `json:"phone" bson:"phone"`
	BusinessPhone string  `json:"businessPhone" bson:"businessPhone"`
}

type To struct {
	Name    string  `json:"name" bson:"name"`
	EmailTo string  `json:"emailTo" bson:"emailTo"`
	Address Address `json:"address" bson:"address"`
	Phone   string  `json:"phone" bson:"phone"`
}

type Address struct {
	Street    string `json:"street" bson:"street"`
	CityState string `json:"cityState" bson:"cityState"`
	ZipCode   string `json:"zipCode" bson:"zipCode"`
}

type InvItem struct {
	Title       string  `json:"title" bson:"title"`
	Description string  `json:"description" bson:"description"`
	Rate        float64 `json:"rate" bson:"rate"`
	Qty         float64 `json:"qty" bson:"qty"`
}
