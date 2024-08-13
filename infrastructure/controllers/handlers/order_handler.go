package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"pm/application"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/persistences/base/logger"
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

// HandleCreateOrder CreateOrder godoc
//
//	@Summary		Create a new order
//	@Description	Create order with order items included
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			CreateOrderRequest	body		payload.CreateOrderRequest	true	"create a new order"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/orders 				[post]
func (h *OrderHandler) HandleCreateOrder(c *gin.Context) {
	ctx, newlogger := logger.GetLogger().Start(c, "handlers/HandleCreateOrder")
	defer newlogger.End()

	var requestPayload payload.CreateOrderRequest
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		newlogger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"data": err.Error()})
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
			newlogger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
			c.Error(payload.ErrInternal(fmt.Errorf("error cast from any to int64")))
			return
		}
	default:
		newlogger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		c.Error(payload.ErrInternal(fmt.Errorf("error cast from any to int64")))
		return
	}

	requestPayload.UserID = uint(id)
	requestPayload.Status = "WAITING_FOR_PAYMENT"
	if err := h.usecase.CreateOrder(ctx, &requestPayload); err != nil {
		newlogger.Error("CREATE_ORDER_FAILED", map[string]interface{}{"message": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
	newlogger.Info("CREATE_ORDER_SUCCESSFULLY", map[string]interface{}{})
}

func (h *OrderHandler) HandleUpdateOrderByID(c *gin.Context) {

}

func (h *OrderHandler) HandleGetAllOrders(c *gin.Context) {

}

// HandleGetOrderByID GetOrderByID godoc
//
//	@Summary		Get order by id
//	@Description	get order by id
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"the id of order to get the order"
//	@Success		200				{object}	payload.AppResponse
//	@Failure		400				{object}	payload.AppError
//	@Failure		404				{object}	payload.AppError
//	@Failure		500				{object}	payload.AppError
//	@Router			/orders/:id 				[get]
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

	orderResponse := mapper.OrderToOrderResponse(order)
	var totalPrice float64 = 0
	for _, v := range order.OrderItems {
		totalPrice += v.Price * float64(v.Quantity)
	}
	orderResponse.Total = math.Round(totalPrice*100) / 100
	utils.HttpSuccessResponse(c, orderResponse, "")
}

func (h *OrderHandler) HandleDeleteOrderByID(c *gin.Context) {

}

func (h *OrderHandler) HandleGetOrdersByUserID(c *gin.Context) {

}
