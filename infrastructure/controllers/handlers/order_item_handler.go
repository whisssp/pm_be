package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"strconv"
)

type OrderItemHandler struct {
	p                *base.Persistence
	orderItemUsecase application.OrderItemUsecase
}

func NewOrderItemHandler(p *base.Persistence) *OrderItemHandler {
	orderItemUsecase := application.NewOrderItemUsecase(p)
	return &OrderItemHandler{p: p, orderItemUsecase: orderItemUsecase}
}

func (oi *OrderItemHandler) HandleCreateNewOrderItem(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleCreateNewOrderItem", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var orderItems []entity.OrderItem
	if err := c.ShouldBindJSON(&orderItems); err != nil {
		oi.p.Logger.Info("BINDING_REQUEST_DATA_TO_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrInvalidRequest(err))
		return
	}

	if err := oi.orderItemUsecase.CreateNewOrderItem(c, orderItems); err != nil {
		oi.p.Logger.Error("CREATING_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("CREATE_ORDER_ITEM: SUCCESSFULLY", map[string]interface{}{"order_items": orderItems})
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

func (oi *OrderItemHandler) HandleGetOrderItemByID(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleGetOrderItemByID", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	id, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		oi.p.Logger.Info("BINDING_ID_PARAM: ERROR", map[string]interface{}{"error": e.Error()})
		c.Error(payload.ErrInvalidRequest(e))
		return
	}

	if id == 0 {
		err := errors.New("param [id] is required")
		oi.p.Logger.Info("INVALID_ORDER_ITEM_ID: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrParamRequired(err))
		return
	}

	orderItem, err := oi.orderItemUsecase.GetOrderItemByID(c, id)
	if err != nil {
		oi.p.Logger.Error("GET_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("GET_ORDER_ITEM: SUCCESSFULLY", map[string]interface{}{"order_item": orderItem})
	c.JSON(http.StatusOK, payload.SuccessResponse(orderItem, ""))
}

func (oi *OrderItemHandler) HandleUpdateOrderItemByID(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleUpdateOrderItemByID", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	id, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		oi.p.Logger.Info("BINDING_ID_PARAM: ERROR", map[string]interface{}{"error": e.Error()})
		c.Error(payload.ErrInvalidRequest(e))
		return
	}

	if id == 0 {
		err := errors.New("param [id] is required")
		oi.p.Logger.Info("INVALID_ORDER_ITEM_ID: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrParamRequired(err))
		return
	}

	var orderItems []entity.OrderItem
	if err := c.ShouldBindJSON(&orderItems); err != nil {
		oi.p.Logger.Info("BINDING_REQUEST_DATA_TO_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrInvalidRequest(err))
		return
	}
	orderItemUpdated, err := oi.orderItemUsecase.UpdateOrderItem(c, orderItems)
	if err != nil {
		oi.p.Logger.Error("UPDATE_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("UPDATE_ORDER_ITEM: SUCCESSFULLY", map[string]interface{}{"order_item": orderItemUpdated})
	c.JSON(http.StatusOK, payload.SuccessResponse(orderItemUpdated, ""))
}

func (oi *OrderItemHandler) HandleDeleteOrderItemByID(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleDeleteOrderItemByID", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	id, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		oi.p.Logger.Info("BINDING_ID_PARAM: ERROR", map[string]interface{}{"error": e.Error()})
		c.Error(payload.ErrInvalidRequest(e))
		return
	}

	if id == 0 {
		err := errors.New("param [id] is required")
		oi.p.Logger.Info("INVALID_ORDER_ITEM_ID: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrParamRequired(err))
		return
	}

	err := oi.orderItemUsecase.DeleteOrderItemByID(c, id)
	if err != nil {
		oi.p.Logger.Error("DELETE_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("DELETE_ORDER_ITEM: SUCCESSFULLY", map[string]interface{}{"id_order_item_deleted": id})
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

func (oi *OrderItemHandler) HandleGetAllOrderItems(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleGetAllOrderItems", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	orderItems, err := oi.orderItemUsecase.GetAllOrderItems(c)
	if err != nil {
		oi.p.Logger.Error("GET_ALL_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("GET_ALL_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": orderItems})
	c.JSON(http.StatusOK, payload.SuccessResponse(orderItems, ""))
}