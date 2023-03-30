package handlers

import (
	"net/http"

	"github.com/avalonprod/gasstrem/src/internal/models"
	"github.com/avalonprod/gasstrem/src/internal/services"
	"github.com/gin-gonic/gin"
)

type InvoiceInput struct {
	UserID      string           `json:"userID" bson:"userID"`
	InvTitle    string           `json:"invTitle" bson:"invTitle"`
	InvNum      int              `json:"invNum" bson:"invNum"`
	CreatedTime string           `json:"createdTime" bson:"createdTime"`
	Balance     float64          `json:"balance" bson:"balance"`
	Notes       string           `json:"notes" bson:"notes"`
	Dispatch    bool             `json:"dispatch" bson:"dispatch"`
	Discount    bool             `json:"discount" bson:"discount"`
	ColorLine   string           `json:"colorLine" bson:"colorLine"`
	Currency    string           `json:"currency" bson:"currency"`
	From        from             `json:"from" bson:"from"`
	To          to               `json:"to" bson:"to"`
	InvList     []models.InvItem `json:"invList" bson:"invList"`
}

type from struct {
	Name          string  `json:"name" bson:"name"`
	EmailFrom     string  `json:"emailFrom" bson:"emailFrom"`
	Address       address `json:"address" bson:"address"`
	Phone         string  `json:"phone" bson:"phone"`
	BusinessPhone string  `json:"businessPhone" bson:"businessPhone"`
}

type to struct {
	Name    string  `json:"name" bson:"name"`
	EmailTo string  `json:"emailTo" bson:"emailTo"`
	Address address `json:"address" bson:"address"`
	Phone   string  `json:"phone" bson:"phone"`
}

type address struct {
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

func (h *Handlers) CreateInvoice(c *gin.Context) {
	var input InvoiceInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	from := services.From{
		Name:          input.From.Name,
		EmailFrom:     input.From.EmailFrom,
		Address:       services.Address(input.From.Address),
		Phone:         input.From.Phone,
		BusinessPhone: input.From.BusinessPhone,
	}

	to := services.To{
		Name:    input.To.Name,
		EmailTo: input.To.EmailTo,
		Address: services.Address(input.To.Address),
		Phone:   input.To.Phone,
	}

	err = h.services.Invoices.CreateInvoice(c.Request.Context(), services.InvoiceInput{
		UserID:      id.String(),
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

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	}
	c.Status(http.StatusCreated)
}
