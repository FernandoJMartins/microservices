package domain

import "time"

type ShippingItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int32  `json:"quantity"`
}

type Shipping struct {
	ID           int64          `json:"id"`
	OrderID      int64          `json:"order_id"`
	Status       string         `json:"status"`
	DeliveryDays int32          `json:"delivery_days"`
	Items        []ShippingItem `json:"items"`
	CreatedAt    int64          `json:"created_at"`
}

func NewShipping(orderId int64, items []ShippingItem) Shipping {
	var totalQuantity int32
	for _, item := range items {
		totalQuantity += item.Quantity
	}
	deliveryDays := int32(1 + totalQuantity/5)

	return Shipping{
		CreatedAt:    time.Now().Unix(),
		Status:       "Pending",
		OrderID:      orderId,
		DeliveryDays: deliveryDays,
		Items:        items,
	}
}
