package repositories

import (
	"database/sql"
	"fmt"
	"product/common"
	"product/datamodels"
)

// OrderRepository 定义订单仓储层需要实现的能力。
type OrderRepository interface {
	Conn() error
	Insert(order *datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

// NewOrderManagerRepository 创建订单仓储层实现。
func NewOrderManagerRepository(table string, sql *sql.DB) OrderRepository {
	return &OrderManagerRepository{table: table, mysqlConn: sql}
}

// OrderManagerRepository 负责直接读写订单表。
type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

// Conn 确保订单仓储层拥有可用连接，并补齐默认表名。
func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.GetMySQLConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

// Insert 向订单表插入订单记录。
func (o *OrderManagerRepository) Insert(order *datamodels.Order) (int64, error) {
	if err := o.Conn(); err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO `%s` (userID,productID,orderStatus) VALUES (?,?,?)", o.table)

	result, err := o.mysqlConn.Exec(
		query,
		order.UserID,
		order.ProductId,
		order.OrderStatus,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Delete 按订单 ID 删除记录。
func (o *OrderManagerRepository) Delete(orderID int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	query := fmt.Sprintf("DELETE FROM `%s` WHERE ID = ?", o.table)
	result, err := o.mysqlConn.Exec(query, orderID)
	if err != nil {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

// Update 按订单 ID 更新所属用户、商品和状态。
func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE `%s` SET userID = ?, productID = ?, orderStatus = ? WHERE ID = ?", o.table)
	_, err := o.mysqlConn.Exec(
		query,
		order.UserID,
		order.ProductId,
		order.OrderStatus,
		order.ID,
	)

	return err
}

// SelectByKey 按订单主键查询单条订单。
func (o *OrderManagerRepository) SelectByKey(orderID int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE ID = ?", o.table)
	rows, err := o.mysqlConn.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := common.RowsToMap(rows)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}
	return common.MapToOrder(result[0]), nil
}

// SelectAll 查询全部订单基础信息。
func (o *OrderManagerRepository) SelectAll() ([]*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM `%s`", o.table)
	rows, err := o.mysqlConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := common.RowsToMap(rows)
	if err != nil {
		return nil, err
	}
	orders := make([]*datamodels.Order, 0, len(result))
	for _, order := range result {
		orders = append(orders, common.MapToOrder(order))
	}
	return orders, nil
}

// SelectAllWithInfo 查询订单及其关联商品信息。
// 返回值使用 map[int]map[string]string，便于模板按字段名直接读取展示。
func (o *OrderManagerRepository) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT o.ID, o.userID, o.productID, o.orderStatus, p.productName, p.productImage, p.productUrl "+
			"FROM `%s` o LEFT JOIN product p ON o.productID = p.ID",
		o.table,
	)

	rows, err := o.mysqlConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := common.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	orderInfo := make(map[int]map[string]string, len(result))
	for _, item := range result {
		order := common.MapToOrder(item)
		orderInfo[int(order.ID)] = map[string]string{
			"ID":           valueToString(item["ID"]),
			"userID":       valueToString(item["userID"]),
			"productID":    valueToString(item["productID"]),
			"orderStatus":  valueToString(item["orderStatus"]),
			"productName":  valueToString(item["productName"]),
			"productImage": valueToString(item["productImage"]),
			"productUrl":   valueToString(item["productUrl"]),
		}
	}

	return orderInfo, nil
}

// valueToString 将任意值安全转换为字符串，便于模板展示。
func valueToString(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
