package db

import (
	"context"
	"fmt"

	"github.com/FernandoJMartins/microservices/shipping/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shipping struct {
	gorm.Model
	OrderID      int64
	Status       string
	DeliveryDays int32
	Items        []ShippingItem
}

type ShippingItem struct {
	gorm.Model
	ProductCode string
	Quantity    int32
	ShippingID  uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	err := db.AutoMigrate(&Shipping{}, &ShippingItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Get(ctx context.Context, id string) (domain.Shipping, error) {
	var shippingEntity Shipping
	res := a.db.WithContext(ctx).Preload("Items").First(&shippingEntity, id)
	var items []domain.ShippingItem
	for _, item := range shippingEntity.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	return domain.Shipping{
		ID:           int64(shippingEntity.ID),
		OrderID:      shippingEntity.OrderID,
		Status:       shippingEntity.Status,
		DeliveryDays: shippingEntity.DeliveryDays,
		Items:        items,
		CreatedAt:    shippingEntity.CreatedAt.UnixNano(),
	}, res.Error
}

func (a Adapter) Save(ctx context.Context, shipping *domain.Shipping) error {
	var items []ShippingItem
	for _, item := range shipping.Items {
		items = append(items, ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	shippingEntity := Shipping{
		OrderID:      shipping.OrderID,
		Status:       shipping.Status,
		DeliveryDays: shipping.DeliveryDays,
		Items:        items,
	}
	res := a.db.WithContext(ctx).Create(&shippingEntity)
	if res.Error == nil {
		shipping.ID = int64(shippingEntity.ID)
	}
	return res.Error
}
