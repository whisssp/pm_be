package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"pm/application"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
)

const (
	userContextKey = "user"
)

type OrderHandler struct {
	p       *base.Persistence
	usecase application.OrderUsecase
}

func NewOrderHandler(p *base.Persistence) *OrderHandler {
	usecase := application.NewOrderUsecase(p)
	return &OrderHandler{p, usecase}
}

func (h *OrderHandler) HandleCreateOrder(c *gin.Context) {
	var requestPayload payload.CreateOrderRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.Error(err)
		return
	}

	//userId := c.Value(userContextKey)
	//if exists == false {
	//	c.Error(payload.ErrInvalidRequest(errors.New("unauthorized because cannot find user from context")))
	//}
	idValue := c.Value(userContextKey)
	var id int64
	var err error

	switch v := idValue.(type) {
	case int64:
		id = v
	case int:
		id = int64(v)
	case int32:
		id = int64(v)
	case float64:
		id = int64(v)
	case string:
		id, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			utils.HttpErrorResponse(c, fmt.Errorf("error parsing string to int64: %v", err))
			return
		}
	default:
		utils.HttpErrorResponse(c, fmt.Errorf("error cast from any to int64"))
		return
	}

	requestPayload.UserID = uint(id)
	requestPayload.Status = "WAITING_FOR_PAYMENT"
	if err := h.usecase.CreateOrder(&requestPayload); err != nil {
		c.Error(err)
		return
	}
	utils.HttpSuccessResponse(c, nil, "")
}

func (h *OrderHandler) HandleUpdateOrderByID(c *gin.Context) {

}

func (h *OrderHandler) HandleGetAllOrders(c *gin.Context) {

}

func (h *OrderHandler) HandleGetOrderByID(c *gin.Context) {
	orderId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if orderId == 0 {
		c.Error(payload.ErrParamRequired(errors.New("param [id] is required")))
		return
	}
	order, err := h.usecase.GetOrderByID(orderId)
	if err != nil {
		c.Error(err)
		return
	}
	utils.HttpSuccessResponse(c, order, "")
}

func (h *OrderHandler) HandleDeleteOrderByID(c *gin.Context) {

}

func (h *OrderHandler) HandleGetOrdersByUserID(c *gin.Context) {

}