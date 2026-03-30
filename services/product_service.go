package services

import (
	"product/datamodels"
	"product/repositories"
)

// IProductService 定义商品业务层对外能力。
type IProductService interface {
	GetProductById(int64) (*datamodels.Product, error)
	GetProductAll() ([]*datamodels.Product, error)
	DeleteProductById(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
}

// ProductService 是商品业务层实现。
// 当前逻辑较薄，主要是对 repository 做一层封装。
type ProductService struct {
	productRepository repositories.IProduct
}

// NewProductService 创建商品业务层实例。
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

// GetProductById 按商品 ID 查询单个商品。
func (p *ProductService) GetProductById(id int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(id)
}

// GetProductAll 查询全部商品。
func (p *ProductService) GetProductAll() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

// DeleteProductById 按 ID 删除商品。
func (p *ProductService) DeleteProductById(id int64) bool {
	return p.productRepository.Delete(id)
}

// InsertProduct 新增商品。
func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

// UpdateProduct 更新商品信息。
func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}
