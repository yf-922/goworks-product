package services

import (
	"encoding/json"
	"product/datamodels"
	"product/repositories"

	"github.com/mediocregopher/radix/v3"
)

const productAllCacheKey = "product:all"

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
	redisPool         *radix.Pool
}

// NewProductService 创建商品业务层实例。
func NewProductService(repository repositories.IProduct, redisPool *radix.Pool) IProductService {
	return &ProductService{
		productRepository: repository,
		redisPool:         redisPool,
	}
}

// GetProductById 按商品 ID 查询单个商品。
func (p *ProductService) GetProductById(id int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(id)
}

// GetProductAll 查询全部商品。
func (p *ProductService) GetProductAll() ([]*datamodels.Product, error) {
	var cacheValue string
	err := p.redisPool.Do(radix.Cmd(&cacheValue, "GET",
		productAllCacheKey))
	if err == nil && cacheValue != "" {
		var products []*datamodels.Product
		if unmarshalErr := json.Unmarshal([]byte(cacheValue), &products); unmarshalErr == nil {
			return products, nil
		}
	}

	products, err := p.productRepository.SelectAll()
	if err != nil {
		return nil, err
	}

	cacheBytes, err := json.Marshal(products)
	if err == nil {
		_ = p.redisPool.Do(radix.Cmd(nil, "SET", productAllCacheKey,
			string(cacheBytes)))
	}

	return products, nil

}

// DeleteProductById 按 ID 删除商品。
func (p *ProductService) DeleteProductById(id int64) bool {

	return p.productRepository.Delete(id)

}

// InsertProduct 新增商品。
func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	productID, err := p.productRepository.Insert(product)
	if err != nil {
		return 0, err
	}
	_ = p.redisPool.Do(radix.Cmd(nil, "DEL", productAllCacheKey))
	return productID, nil
}

// UpdateProduct 更新商品信息。
func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	if err := p.productRepository.Update(product); err != nil {
		return err
	}
	_ = p.redisPool.Do(radix.Cmd(nil, "DEL", productAllCacheKey))
	return nil
}
