package services

import (
	"context"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"github.com/avalonprod/gasstrem/src/internal/storages"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoicesService struct {
	storages storages.Invoices
}

func NewInvoicesService(storages storages.Invoices) *InvoicesService {
	return &InvoicesService{
		storages: storages,
	}
}

type InvoiceInput struct {
	UserID      primitive.ObjectID `json:"userID" bson:"userID"`
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
	InvList     []models.InvItem   `json:"invList" bson:"invList"`
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

type invItem struct {
	Title       string  `json:"title" bson:"title"`
	Description string  `json:"description" bson:"description"`
	Rate        float64 `json:"rate" bson:"rate"`
	Qty         float64 `json:"qty" bson:"qty"`
}

func (i *InvoicesService) CreateInvoice(ctx context.Context, input InvoiceInput) error {

	from := models.From{
		Name:          input.From.Name,
		EmailFrom:     input.From.EmailFrom,
		Address:       models.Address(input.From.Address),
		Phone:         input.From.Phone,
		BusinessPhone: input.From.BusinessPhone,
	}

	to := models.To{
		Name:    input.To.Name,
		EmailTo: input.To.EmailTo,
		Address: models.Address(input.To.Address),
		Phone:   input.To.Phone,
	}

	err := i.storages.Create(ctx, models.Invoice{
		UserID:      input.UserID,
		InvTitle:    input.InvTitle,
		InvNum:      input.InvNum,
		CreatedTime: input.CreatedTime,
		Balance:     input.Balance,
		Notes:       input.Notes,
		Dispatch:    input.Dispatch,
		Discount:    input.Discount,
		ColorLine:   input.ColorLine,
		Currency:    input.Currency,
		From:        from,
		To:          to,
		InvList:     input.InvList,
	})

	return err
}

func (i *InvoicesService) GetAllInvoceByUserId(ctx context.Context, userID primitive.ObjectID) ([]models.Invoice, error) {
	res, err := i.storages.GetAllInvoceByUserId(ctx, userID)

	if err != nil {
		return []models.Invoice{}, err
	}

	return res, nil
}

func GenerateInvoceTemplate() {

}
