package application

import (
	"github.com/gin-gonic/gin"
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
	orderItems "pm/infrastructure/implementations/order_items"
	"pm/infrastructure/mapper"
	"pm/infrastructure/persistences/base"
	"pm/utils"
)

type OrderItemUsecase interface {
	CreateNewOrderItem(*gin.Context, []entity.OrderItem) error
	UpdateOrderItem(*gin.Context, []entity.OrderItem) ([]payload.OrderItemResponse, error)
	GetOrderItemByID(*gin.Context, int64) (*payload.OrderItemResponse, error)
	DeleteOrderItemByID(*gin.Context, int64) error
	GetAllOrderItems(*gin.Context) (*payload.ListOrderItemResponses, error)
}

type orderItemUsecase struct {
	p *base.Persistence
}

// CreateNewOrderItem HandleCreateNewOrderItem godoc
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
func (o orderItemUsecase) CreateNewOrderItem(c *gin.Context, items []entity.OrderItem) error {
	span := o.p.Logger.Start(c, "CREATE_ORDER_ITEM_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	if err := utils.ValidateReqPayload(items); err != nil {
		o.p.Logger.Error("CREATE_ORDER_ITEMS: ERROR INVALID DATA", map[string]interface{}{"error": err.Error(), "order_items": items})
		return payload.ErrInvalidRequest(err)
	}

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	if err := oiRepo.CreateNewOrderItems(items); err != nil {
		o.p.Logger.Error("CREATE_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
		return err
	}

	o.p.Logger.Info("CREATE_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": items})
	return nil
}

// UpdateOrderItem HandleUpdateOrderItemByID godoc
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
func (o orderItemUsecase) UpdateOrderItem(c *gin.Context, items []entity.OrderItem) ([]payload.OrderItemResponse, error) {
	span := o.p.Logger.Start(c, "UPDATEE_ORDER_ITEM_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	if err := utils.ValidateReqPayload(items); err != nil {
		o.p.Logger.Error("CREATE_ORDER_ITEMS: ERROR INVALID DATA", map[string]interface{}{"error": err.Error(), "order_items": items})
		return nil, payload.ErrInvalidRequest(err)
	}

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	orderItems, err := oiRepo.UpdateOrderItems(items)
	if err != nil {
		o.p.Logger.Error("UPDATE_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
		return nil, err
	}

	orderItemResponses := make([]payload.OrderItemResponse, 0)
	for _, item := range orderItems {
		oir := mapper.OrderItemToOrderItemResponse(&item)
		orderItemResponses = append(orderItemResponses, oir)
	}

	o.p.Logger.Info("UPDATE_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": items})
	return orderItemResponses, nil
}

// GetOrderItemByID HandleGetOrderItemByID godoc
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
func (o orderItemUsecase) GetOrderItemByID(c *gin.Context, id int64) (*payload.OrderItemResponse, error) {
	span := o.p.Logger.Start(c, "GET_ORDER_ITEM_BY_ID_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	orderItem, err := oiRepo.GetOrderItemByID(id)
	if err != nil {
		o.p.Logger.Error("GET_ORDER_ITEM_BY_ID: ERROR", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	oir := mapper.OrderItemToOrderItemResponse(orderItem)

	o.p.Logger.Info("GET_ORDER_ITEM_BY_ID: SUCCESSFULLY", map[string]interface{}{"order_item": oir})
	return &oir, nil
}

// DeleteOrderItemByID HandleDeleteOrderItemByID godoc
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
func (o orderItemUsecase) DeleteOrderItemByID(c *gin.Context, id int64) error {
	span := o.p.Logger.Start(c, "DELETE_ORDER_ITEM_BY_ID_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	err := oiRepo.DeleteOrderItemByID(id)
	if err != nil {
		o.p.Logger.Error("DELETE_ORDER_ITEM_BY_ID: ERROR", map[string]interface{}{"error": err.Error()})
		return err
	}

	o.p.Logger.Info("DELETE_ORDER_ITEM_BY_ID: SUCCESSFULLY", map[string]interface{}{"deleted_id": id})
	return nil
}

// GetAllOrderItems HandleGetAllOrderItems godoc
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
func (o orderItemUsecase) GetAllOrderItems(c *gin.Context) (*payload.ListOrderItemResponses, error) {
	span := o.p.Logger.Start(c, "GET_ALL_ORDER_ITEMS_USECASE", o.p.Logger.SetContextWithSpanFunc())
	defer span.End()

	oiRepo := orderItems.NewOrderItemRepository(o.p.GormDB, c, o.p, span)
	items, err := oiRepo.GetAllOrderItems()
	if err != nil {
		o.p.Logger.Error("GET_ALL_ORDER_ITEMS: ERROR", map[string]interface{}{"error": err.Error(), "order_items": items})
		return nil, err
	}

	orderItemResponses := make([]payload.OrderItemResponse, 0)
	for _, item := range items {
		oir := mapper.OrderItemToOrderItemResponse(&item)
		orderItemResponses = append(orderItemResponses, oir)
	}
	response := payload.ListOrderItemResponses{Orders: orderItemResponses}

	o.p.Logger.Info("GET_ALL_ORDER_ITEMS: SUCCESSFULLY", map[string]interface{}{"order_items": response})
	return &response, nil
}

func NewOrderItemUsecase(p *base.Persistence) OrderItemUsecase {
	return orderItemUsecase{p: p}
}