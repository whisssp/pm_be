package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pm/application"
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

// HandleCreateNewOrderItem godoc
//
//	@Summary		Create new order item
//	@Description	Create a new order item
//	@Tags			OrderItem
//	@Accept			json
//	@Produce		json
//	@Param			orderItems	body		[]entity.OrderItem	true	"Order items to create"
//	@Success		200			{object}	payload.AppResponse
//	@Failure		400			{object}	payload.AppError
//	@Failure		500			{object}	payload.AppError
//	@Router			/order-items	[post]
func (oi *OrderItemHandler) HandleCreateNewOrderItem(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleCreateNewOrderItem", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	var orderItemRequest payload.OrderItemsRequest
	if err := c.ShouldBindJSON(&orderItemRequest); err != nil {
		oi.p.Logger.Info("BINDING_REQUEST_DATA_TO_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrInvalidRequest(err))
		return
	}

	if err := oi.orderItemUsecase.CreateNewOrderItem(c, orderItemRequest.OrderItems); err != nil {
		oi.p.Logger.Error("CREATING_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("CREATE_ORDER_ITEM: SUCCESSFULLY", map[string]interface{}{"order_items": orderItemRequest.OrderItems})
	c.JSON(http.StatusOK, payload.SuccessResponse(nil, ""))
}

// HandleGetOrderItemByID godoc
//
//	@Summary		Get order item by ID
//	@Description	Get an order item by its ID
//	@Tags			OrderItem
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID of the order item"
//	@Success		200	{object}	payload.AppResponse
//	@Failure		400	{object}	payload.AppError
//	@Failure		404	{object}	payload.AppError
//	@Failure		500	{object}	payload.AppError
//	@Router			/order-items/{id}	[get]
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

// HandleUpdateOrderItemByID godoc
//
//	@Summary		Update order item by ID
//	@Description	Update an order item by its ID
//	@Tags			OrderItem
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"ID of the order item"
//	@Param			orderItem	body		entity.OrderItem	true	"Order item data to update"
//	@Success		200			{object}	payload.AppResponse
//	@Failure		400			{object}	payload.AppError
//	@Failure		404			{object}	payload.AppError
//	@Failure		500			{object}	payload.AppError
//	@Router			/order-items/{id}	[put]
func (oi *OrderItemHandler) HandleUpdateOrderItemByID(c *gin.Context) {
	span := oi.p.Logger.Start(c, "handlers/HandleUpdateOrderItemByID", oi.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	//id, e := strconv.ParseInt(c.Param("id"), 10, 64)
	//if e != nil {
	//	oi.p.Logger.Info("BINDING_ID_PARAM: ERROR", map[string]interface{}{"error": e.Error()})
	//	c.Error(payload.ErrInvalidRequest(e))
	//	return
	//}
	//
	//if id == 0 {
	//	err := errors.New("param [id] is required")
	//	oi.p.Logger.Info("INVALID_ORDER_ITEM_ID: ERROR", map[string]interface{}{"error": err.Error()})
	//	c.Error(payload.ErrParamRequired(err))
	//	return
	//}

	var orderItemRequest payload.OrderItemsRequest
	if err := c.ShouldBindJSON(&orderItemRequest); err != nil {
		oi.p.Logger.Info("BINDING_REQUEST_DATA_TO_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(payload.ErrInvalidRequest(err))
		return
	}
	orderItemUpdated, err := oi.orderItemUsecase.UpdateOrderItem(c, orderItemRequest.OrderItems)
	if err != nil {
		oi.p.Logger.Error("UPDATE_ORDER_ITEM: ERROR", map[string]interface{}{"error": err.Error()})
		c.Error(err)
		return
	}

	oi.p.Logger.Info("UPDATE_ORDER_ITEM: SUCCESSFULLY", map[string]interface{}{"order_item": orderItemUpdated})
	c.JSON(http.StatusOK, payload.SuccessResponse(orderItemUpdated, ""))
}

// HandleDeleteOrderItemByID godoc
//
//	@Summary		Delete order item by ID
//	@Description	Delete an order item by its ID
//	@Tags			OrderItem
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID of the order item to delete"
//	@Success		200	{object}	payload.AppResponse
//	@Failure		400	{object}	payload.AppError
//	@Failure		404	{object}	payload.AppError
//	@Failure		500	{object}	payload.AppError
//	@Router			/order-items/{id}	[delete]
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

// HandleGetAllOrderItems godoc
//
//	@Summary		Get all order items
//	@Description	Get a list of all order items
//	@Tags			OrderItem
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	payload.AppResponse
//	@Failure		400	{object}	payload.AppError
//	@Failure		500	{object}	payload.AppError
//	@Router			/order-items	[get]
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