package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	span := h.p.Logger.Start(c, "handlers/HandleCreateOrder", h.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var requestPayload payload.CreateOrderRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		h.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"data": err.Error()})
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
			h.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
			c.Error(payload.ErrInternal(fmt.Errorf("error cast from any to int64")))
			return
		}
	default:
		h.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		c.Error(payload.ErrInternal(fmt.Errorf("error cast from any to int64")))
		return
	}

	requestPayload.UserID = uint(id)
	requestPayload.Status = "WAITING_FOR_PAYMENT"
	if err := h.usecase.CreateOrder(c, &requestPayload); err != nil {
		h.p.Logger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
	h.p.Logger.Info("CREATE_ORDER_SUCCESSFULLY", map[string]interface{}{})
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
	order, err := h.usecase.GetOrderByID(c, orderId)
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