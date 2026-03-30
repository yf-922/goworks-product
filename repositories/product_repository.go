package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"product/common"
	"product/datamodels"
)

// IProduct 定义商品仓储层对外提供的能力。
// service 层只依赖接口，不直接依赖具体实现。
type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

// ProductManager 是商品仓储层的默认实现。
// 它负责直接操作 product 表。
type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

// NewProductManager 创建商品仓储层实例。
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

// Conn 确保仓储层持有可用数据库连接。
// 如果外部没有注入连接，这里会按默认配置自行创建。
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.GetMySQLConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	return nil
}

// Insert 向商品表插入一条新记录，并返回数据库生成的主键。
func (p *ProductManager) Insert(product *datamodels.Product) (int64, error) {
	if err := p.Conn(); err != nil {
		return 0, err
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (productName, productNum, productImage, productUrl) VALUES (?, ?, ?, ?)",
		p.table,
	)

	result, err := p.mysqlConn.Exec(
		query,
		product.ProductName,
		product.ProductNum,
		product.ProductImage,
		product.ProductUrl,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Delete 按商品 ID 删除记录。
// 返回值表示是否至少删除了一条数据。
func (p *ProductManager) Delete(productID int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE ID = ?", p.table)
	result, err := p.mysqlConn.Exec(query, productID)
	if err != nil {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

// Update 按商品 ID 更新商品的主要字段。
func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	query := fmt.Sprintf(
		"UPDATE %s SET productName = ?, productNum = ?, productImage = ?, productUrl = ? WHERE ID = ?",
		p.table,
	)

	result, err := p.mysqlConn.Exec(
		query,
		product.ProductName,
		product.ProductNum,
		product.ProductImage,
		product.ProductUrl,
		product.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

// SelectByKey 按商品主键查询单条记录。
func (p *ProductManager) SelectByKey(productID int64) (*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID = ?", p.table)
	rows, err := p.mysqlConn.Query(query, productID)
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

	return common.MapToProduct(result[0]), nil
}

// SelectAll 查询整张商品表，并逐条映射为 Product 结构体切片。
func (p *ProductManager) SelectAll() ([]*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM %s", p.table)
	rows, err := p.mysqlConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := common.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	products := make([]*datamodels.Product, 0, len(result))
	for _, item := range result {
		products = append(products, common.MapToProduct(item))
	}
	return products, nil
}
