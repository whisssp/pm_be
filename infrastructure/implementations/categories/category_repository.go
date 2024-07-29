package categories

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"pm/domain/entity"
	"pm/domain/repository"
	"pm/infrastructure/controllers/payload"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return CategoryRepository{db}
}

func (c CategoryRepository) Create(category *entity.Category) error {
	db := c.db
	if err := db.Create(category).Error; err != nil {
		return payload.ErrDB(err)
	}
	return nil
}

func (c CategoryRepository) Update(category *entity.Category) (*entity.Category, error) {
	db := c.db
	if err := db.Updates(category).Error; err != nil {
		return nil, payload.ErrDB(err)
	}
	return category, nil
}

func (c CategoryRepository) GetCategoryByID(id int64) (*entity.Category, error) {
	var category entity.Category
	db := c.db
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("categories", err)
		}
		return nil, payload.ErrDB(err)
	}
	return &category, nil
}

func (c CategoryRepository) GetAllCategories(filter *entity.CategoryFilter, pagination *entity.Pagination) ([]entity.Category, error) {
	categories := make([]entity.Category, 0)
	var totalRows int64
	db := c.db
	db = db.Model(&entity.Category{}).Count(&totalRows)
	if filter != nil {
		db = db.Scopes(applyFilter(filter)).Count(&totalRows)
	}
	if err := db.Scopes(paginate(pagination)).Find(&categories).Error; err != nil {
		return nil, payload.ErrDB(err)
	}
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	return categories, nil
}

func (c CategoryRepository) DeleteCategory(category *entity.Category) error {
	db := c.db
	if err := db.Delete(category).Error; err != nil {
		return payload.ErrDB(err)
	}
	return nil
}

func paginate(pagination *entity.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func applyFilter(filter *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(
			applyFilterKeyword(filter),
			applyFilterID(filter),
			applyFilterName(filter),
			applyCreatedAtFilter(filter),
			applyUpdatedAtFilter(filter),
			applyDeletedFilter(filter),
		)
	}
}

func applyFilterID(f *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.ID != 0 {
			return db.Where("id = ?", f.ID)
		}
		return db
	}
}

func applyFilterKeyword(f *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Keyword != "" {
			return db.Where("name ILIKE ? OR CAST(id as text) LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
		}
		return db
	}
}

func applyFilterName(f *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Name != "" {
			return db.Where("cast(unaccent(name) as text) ILIKE ?", "%"+f.Keyword+"%")
		}
		return db
	}
}

func applyCreatedAtFilter(f *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
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

func applyUpdatedAtFilter(f *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
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

func applyDeletedFilter(f *entity.CategoryFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.Deleted {
			db = db.Where("deleted = ?", f.Deleted)
		}
		return db
	}
}