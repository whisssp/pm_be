package categories

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

type CategoryRepository struct {
	db *gorm.DB
	p  *base.Persistence
	c  *gin.Context
}

func NewCategoryRepository(c *gin.Context, p *base.Persistence, db *gorm.DB) repository.CategoryRepository {
	return CategoryRepository{db, p, c}
}

func (c CategoryRepository) Create(category *entity.Category) error {
	span := c.p.Logger.Start(c.c, "CREATE_CATEGORY_DATABASE")
	defer span.End()
	c.p.Logger.Info("CREATE_CATEGORY", map[string]interface{}{"data": category})

	db := c.db
	if err := db.Create(&category).Error; err != nil {
		c.p.Logger.Error("CREATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return payload.ErrDB(err)
	}

	c.p.Logger.Info("CREATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": category})
	return nil
}

func (c CategoryRepository) Update(category *entity.Category) (*entity.Category, error) {
	span := c.p.Logger.Start(c.c, "UPDATE_CATEGORY_DATABASE")
	defer span.End()
	c.p.Logger.Info("UPDATE_CATEGORY", map[string]interface{}{"data": category})

	db := c.db
	if err := db.Debug().Model(&category).Updates(category).Error; err != nil {
		c.p.Logger.Info("UPDATE_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, payload.ErrEntityNotFound("products", err)

		}
		return nil, payload.ErrDB(err)
	}

	c.p.Logger.Info("UPDATE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": category})
	return category, nil
}

func (c CategoryRepository) GetCategoryByID(id int64) (*entity.Category, error) {
	span := c.p.Logger.Start(c.c, "GET_CATEGORY_DATABASE")
	defer span.End()
	c.p.Logger.Info("GET_CATEGORY", map[string]interface{}{"data": id})

	var category entity.Category
	db := c.db
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.p.Logger.Error("GET_CATEGORY_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, payload.ErrDB(err)
	}

	c.p.Logger.Info("GET_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": category})
	return &category, nil
}

func (c CategoryRepository) GetAllCategories(filter *entity.CategoryFilter, pagination *entity.Pagination) ([]entity.Category, error) {
	span := c.p.Logger.Start(c.c, "GET_ALL_CATEGORIES_DATABASE")
	defer span.End()
	c.p.Logger.Info("GET_ALL_CATEGORIES", map[string]interface{}{
		"params": struct {
			Filter     *entity.CategoryFilter `json:"filter"`
			Pagination *entity.Pagination     `json:"pagination"`
		}{
			Filter:     filter,
			Pagination: pagination,
		},
	})

	categories := make([]entity.Category, 0)
	var totalRows int64
	db := c.db
	db = db.Model(&entity.Category{}).Debug().Count(&totalRows)
	if filter != nil {
		db = db.Scopes(applyFilter(filter)).Count(&totalRows)
	}
	if err := db.Scopes(paginate(pagination)).Find(&categories).Error; err != nil {
		c.p.Logger.Error("GET_ALL_CATEGORIES_FAILED", map[string]interface{}{"message": err.Error()})
		return nil, payload.ErrDB(err)
	}
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	c.p.Logger.Info("GET_ALL_CATEGORIES_SUCCESSFULLY", map[string]interface{}{"raw_data": categories, "pagination": pagination})
	return categories, nil
}

func (c CategoryRepository) DeleteCategory(category *entity.Category) error {
	span := c.p.Logger.Start(c.c, "DELETE_CATEGORY_DATABASE")
	defer span.End()
	c.p.Logger.Info("DELETE_CATEGORY", map[string]interface{}{"data": category})

	db := c.db
	if err := db.Delete(&category).Error; err != nil {
		c.p.Logger.Error("DELETE_CATEGORY_FAILED", map[string]interface{}{"data": category, "message": err.Error()})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload.ErrEntityNotFound("products", err)

		}
		return payload.ErrDB(err)
	}

	c.p.Logger.Info("DELETE_CATEGORY_SUCCESSFULLY", map[string]interface{}{"data": category})
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