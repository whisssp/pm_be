package products

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductRepository{db}
}

func (prodRepo *ProductRepository) Create(product *entity.Product) error {
	err := prodRepo.db.Create(product).Error
	if err != nil {
		fmt.Printf("error creating product on database: %v", err)
		return payload.ErrDB(err)
	}
	return nil
}

func (prodRepo *ProductRepository) Update(product *entity.Product) (*entity.Product, error) {
	if err := prodRepo.db.Updates(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (prodRepo *ProductRepository) GetProductByID(id int64) (*entity.Product, error) {
	var product entity.Product
	if err := prodRepo.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (prodRepo *ProductRepository) GetAllProducts(filter *payload.ProductFilter, pagination *payload.Pagination) ([]entity.Product, error) {
	var totalRows int64
	products := make([]entity.Product, 0)
	db := prodRepo.db.Model(entity.Product{})
	if filter != nil {
		db = db.Scopes(applyFilter(filter)).Count(&totalRows)
	}
	if pagination != nil {
		if err := db.Scopes(paginate(pagination)).Find(&products).Count(&totalRows).Error; err != nil {
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
	if err := prodRepo.db.Delete(product).Error; err != nil {
		return err
	}
	return nil
}

func paginate(pagination *payload.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func applyFilter(f *payload.ProductFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(
			ApplyKeywordFilter(f, db),
			ApplyIDFilter(f, db),
			ApplyNameFilter(f, db),
			ApplyPriceFilter(f, db),
			ApplyDescriptionFilter(f, db),
			ApplyCategoryIDFilter(f, db),
			ApplyCreatedAtFilter(f, db),
			ApplyUpdatedAtFilter(f, db),
			ApplyDeletedFilter(f, db),
		)
	}
}

func ApplyKeywordFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Keyword != "" {
			db = db.Where("name LIKE ? OR description LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
		}
		return db
	}
}

func ApplyIDFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.ID != 0 {
			db = db.Where("id = ?", f.ID)
		}
		return db
	}
}

func ApplyNameFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Name != "" {
			db = db.Where("name LIKE ?", "%"+f.Name+"%")
		}
		return db
	}
}

func ApplyPriceFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
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

func ApplyDescriptionFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Description != "" {
			db = db.Where("description LIKE ?", "%"+f.Description+"%")
		}
		return db
	}
}

func ApplyCategoryIDFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.CategoryID != 0 {
			db = db.Where("category_id = ?", f.CategoryID)
		}
		return db
	}
}

func ApplyCreatedAtFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
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

func ApplyUpdatedAtFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
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

func ApplyDeletedFilter(f *payload.ProductFilter, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Deleted {
			db = db.Where("deleted = ?", f.Deleted)
		}
		return db
	}
}