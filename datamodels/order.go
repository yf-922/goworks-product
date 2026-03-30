package datamodels

// Order 对应订单表的数据结构。
type Order struct {
	// ID 是订单主键。
	ID int64 `sql:"ID"`
	// UserID 是下单用户 ID。
	UserID int64 `sql:"userID"`
	// ProductId 是关联商品 ID。
	ProductId int64 `sql:"productId"`
	// OrderStatus 是订单状态码。
	OrderStatus int64 `sql:"orderStatus"`
}

const (
	// OrderWait 表示订单待处理。
	OrderWait = iota
	// OrederSuccess 表示订单成功。
	// 常量名存在拼写问题，但当前项目里已经在使用，暂不改动行为。
	OrederSuccess
	// OrderFailed 表示订单失败。
	OrderFailed
)
