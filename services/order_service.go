package services

import (
	"database/sql"
	"errors"
	"product/datamodels"
	"product/repositories"
)

var ErrOrderUserNotFound = errors.New("user not found")
var ErrOrderProductNotFound = errors.New("product not found")

// IOrderService 定义订单业务层对外提供的方法。
type IOrderService interface {
	GetOrderById(int64) (*datamodels.Order, error)
	DeleteOrderById(int64) bool
	UpdateOrderById(order *datamodels.Order) error
	InsertOrder(order *datamodels.Order) (int64, error)
	GetAllOrders() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
}

// NewOrderService 创建订单业务层实例。
func NewOrderService(orderRepository repositories.OrderRepository, productRepository repositories.IProduct,
	userRepository repositories.IUserRepository) IOrderService {
	return &OrderService{
		OrderRepository:   orderRepository,
		ProductRepository: productRepository,
		UserRepository:    userRepository,
	}
}

// OrderService 是订单业务层实现。
type OrderService struct {
	OrderRepository   repositories.OrderRepository
	ProductRepository repositories.IProduct
	UserRepository    repositories.IUserRepository
}

// GetOrderById 按订单 ID 查询单个订单。
func (o *OrderService) GetOrderById(id int64) (*datamodels.Order, error) {
	return o.OrderRepository.SelectByKey(id)
}

// DeleteOrderById 按 ID 删除订单。
func (o *OrderService) DeleteOrderById(id int64) bool {
	return o.OrderRepository.Delete(id)
}

// UpdateOrderById 更新订单内容。
func (o *OrderService) UpdateOrderById(order *datamodels.Order) error {
	return o.OrderRepository.Update(order)
}

// InsertOrder 新增订单。
func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	_, err := o.UserRepository.SelectById(order.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrOrderUserNotFound
		}
		return 0, err
	}

	_, err = o.ProductRepository.SelectByKey(order.ProductId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrOrderProductNotFound
		}
		return 0, err
	}

	return o.OrderRepository.Insert(order)
}

// GetAllOrders 查询全部订单基础信息。
func (o *OrderService) GetAllOrders() ([]*datamodels.Order, error) {
	return o.OrderRepository.SelectAll()
}

// GetAllOrderInfo 查询全部订单及其关联商品展示信息。
func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectAllWithInfo()
}
