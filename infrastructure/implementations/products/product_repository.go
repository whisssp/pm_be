package products

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

const (
	entityName string = "products"
)

type ProductRepository struct {
	db *gorm.DB
	p  *base.Persistence
	c  *gin.Context
}

func NewProductRepository(c *gin.Context, p *base.Persistence, db *gorm.DB) repository.ProductRepository {
	if c == nil {
		return &ProductRepository{
			db: db,
			p:  p,
			c:  nil,
		}
	}
	return &ProductRepository{db, p, c}
}

func (prodRepo *ProductRepository) Create(product *entity.Product) error {
	db := prodRepo.db
	err := db.Create(product).Error
	if err != nil {
		return payload.ErrDB(err)
	}
	return nil
}

func (prodRepo *ProductRepository) Update(product *entity.Product) (*entity.Product, error) {
	db := prodRepo.db
	if err := db.Updates(product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound(entityName, err)
		}
		return nil, payload.ErrDB(err)
	}
	return product, nil
}

func (prodRepo *ProductRepository) UpdateMultiProduct(products ...entity.Product) ([]entity.Product, error) {
	span := prodRepo.p.Logger.Start(prodRepo.c, "UPDATE_PRODUCT_STOCK")
	defer span.End()

	tx := prodRepo.db.Begin()
	prodRepo.p.Logger.Info("UPDATE_PRODUCT_STOCK", map[string]interface{}{"products": products})
	for index, p := range products {
		if err := tx.Model(&products[index]).Update("stock", p.Stock).Error; err != nil {
			prodRepo.p.Logger.Error("UPDATE_PRODUCT_STOCK_FAILED", map[string]interface{}{"product": p, "message": err.Error()})
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, payload.ErrEntityNotFound(entityName, err)
			}
			return nil, payload.ErrDB(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		prodRepo.p.Logger.Error("UPDATE_PRODUCT_STOCK_FAILED", map[string]interface{}{"products": products, "message": err.Error()})
		return nil, payload.ErrDB(err)
	}
	prodRepo.p.Logger.Error("UPDATE_PRODUCT_STOCK_SUCCESSFULLY", map[string]interface{}{"products": products})
	return products, nil
}

func (prodRepo *ProductRepository) GetProductByID(id int64) (*entity.Product, error) {
	db := prodRepo.db
	var product entity.Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, payload.ErrDB(err)
	}
	return &product, nil
}

func (prodRepo *ProductRepository) GetAllProducts(filter *entity.ProductFilter, pagination *entity.Pagination) ([]entity.Product, error) {
	var totalRows int64
	products := make([]entity.Product, 0)
	db := prodRepo.db
	db = db.Model(entity.Product{}).Count(&totalRows)
	if filter != nil {
		db = db.Scopes(applyFilter(filter)).Count(&totalRows)
	}
	if pagination != nil {
		if err := db.Scopes(paginate(pagination)).Find(&products).Error; err != nil {
			return nil, payload.ErrDB(err)
		}
		pagination.TotalRows = totalRows
		totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
		pagination.TotalPages = totalPages
		return products, nil
	}
	if err := db.Find(&products).Error; err != nil {
		return nil, payload.ErrDB(err)
	}
	return products, nil
}

func (prodRepo *ProductRepository) DeleteProduct(product *entity.Product) error {
	db := prodRepo.db
	if err := db.Delete(product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload.ErrEntityNotFound(entityName, err)
		}
		return payload.ErrDB(err)
	}
	return nil
}

func (prodRepo *ProductRepository) GetStockByProductIDs(orderItems ...entity.OrderItem) ([]entity.Product, error) {
	span := prodRepo.p.Logger.Start(prodRepo.c, "GET_PRODUCT_ID_FROM_ORDER_ITEM")
	defer span.End()

	keyArr := make([]uint, len(orderItems))
	for i, _ := range keyArr {
		keyArr[i] = orderItems[i].ProductID
	}
	prodRepo.p.Logger.Info("GET_PRODUCT", map[string]interface{}{"productIds": keyArr})
	products := make([]entity.Product, len(orderItems))
	err := prodRepo.db.Where("id IN ?", keyArr).Find(&products)
	if err.Error != nil {
		errD := err.Error
		prodRepo.p.Logger.Error("GET_PRODUCT_ERROR", map[string]interface{}{"products": products, "message": errD.Error()})
		return nil, errD
	}
	prodRepo.p.Logger.Info("GET_PRODUCT_SUCCESSFULLY", map[string]interface{}{"products": products})
	return products, nil
}

func paginate(pagination *entity.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func applyFilter(f *entity.ProductFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(
			applyKeywordFilter(f, db),
			applyIDFilter(f, db),
			applyNameFilter(f, db),
			applyPriceFilter(f, db),
			applyDescriptionFilter(f, db),
			applyCategoryIDFilter(f, db),
			applyCreatedAtFilter(f, db),
			applyUpdatedAtFilter(f, db),
			applyDeletedFilter(f, db),
		)
	}
}

func applyKeywordFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Keyword != "" {
			db = db.Where("name LIKE ? OR description LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
		}
		return db
	}
}

func applyIDFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.ID != 0 {
			db = db.Where("id = ?", f.ID)
		}
		return db
	}
}

func applyNameFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Name != "" {
			db = db.Where("name LIKE ?", "%"+f.Name+"%")
		}
		return db
	}
}

func applyPriceFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.PriceFrom != 0 {
			db = db.Where("price >= ?", f.PriceFrom)
		}
		if f.PriceTo != 0 {
			db = db.Where("price <= ?", f.PriceTo)
		}
		return db
	}
}

func applyDescriptionFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Description != "" {
			db = db.Where("description LIKE ?", "%"+f.Description+"%")
		}
		return db
	}
}

func applyCategoryIDFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.CategoryID != 0 {
			db = db.Where("category_id = ?", f.CategoryID)
		}
		return db
	}
}

func applyCreatedAtFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.CreatedAtFrom != nil {
			db = db.Where("created_at >= ?", *f.CreatedAtFrom)
		}
		if f.CreatedAtTo != nil {
			db = db.Where("created_at <= ?", *f.CreatedAtTo)
		}
		return db
	}
}

func applyUpdatedAtFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.UpdatedAtFrom != nil {
			db = db.Where("updated_at >= ?", *f.UpdatedAtFrom)
		}
		if f.UpdatedAtTo != nil {
			db = db.Where("updated_at <= ?", *f.UpdatedAtTo)
		}
		return db
	}
}

func applyDeletedFilter(f *entity.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Deleted {
			db = db.Where("deleted = ?", f.Deleted)
		}
		return db
	}
}