package mapper

import (
	"pm/domain/entity"
	"pm/infrastructure/controllers/payload"
)

func CreateOrderPayloadToOrder(reqPayload *payload.CreateOrderRequest) entity.Order {
	orderItems := OrderItemRequestsToOrderItems(reqPayload.OrderItems)
	return entity.Order{
		UserID:     reqPayload.UserID,
		Status:     reqPayload.Status,
		OrderItems: orderItems,
	}
}

func OrderItemRequestsToOrderItems(items []payload.OrderItemRequest) []entity.OrderItem {
	result := make([]entity.OrderItem, 0)
	for _, v := range items {
		orderItem := entity.OrderItem{
			OrderID:   0,
			ProductID: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
		}
		result = append(result, orderItem)
	}
	return result
}

//func OrderItemRequestToOrderItem(payload *payload.OrderItemRequest) interface{} {
//
//}

func OrderToOrderResponse(e *entity.Order) payload.OrderResponse {
	orderItems := OrderItemsToOrderItemResponses(e.OrderItems)
	return payload.OrderResponse{
		ID:         int64(e.ID),
		UserID:     int64(e.UserID),
		Status:     e.Status,
		OrderItems: orderItems,
		Total:      0,
		AuditTime: payload.AuditTime{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func OrderItemToOrderItemResponse(e *entity.OrderItem) payload.OrderItemResponse {
	return payload.OrderItemResponse{
		ID:        int64(e.ID),
		ProductID: int64(e.ProductID),
		Quantity:  e.Quantity,
		Price:     e.Price,
		AuditTime: payload.AuditTime{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func OrderItemsToOrderItemResponses(listEntity []entity.OrderItem) []payload.OrderItemResponse {
	orderItemResponses := make([]payload.OrderItemResponse, 0)
	for _, v := range listEntity {
		itemResponse := OrderItemToOrderItemResponse(&v)
		orderItemResponses = append(orderItemResponses, itemResponse)
	}
	return orderItemResponses
}