package product

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"inventory/internal/entity"
	"inventory/internal/pkg/apperror"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

func (r *productRepository) Create(ctx context.Context, product *entity.Product) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(product).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperror.NewError(40900000, err, "product", product.ProductID)
		}
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *productRepository) CreateBulk(ctx context.Context, products []entity.Product) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(&products).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Save(product).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *productRepository) Patch(ctx context.Context, product *entity.Product) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Model(&entity.Product{}).Where("product_id = ?", product.ProductID).Updates(product).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *productRepository) Find(ctx context.Context, filter *entity.Product) (*entity.Product, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	primaryFilter, appErr := r.buildPrimaryFilter(filter)
	if appErr != nil {
		return nil, appErr
	}

	var product entity.Product
	if err := db.WithContext(ctx).Where(primaryFilter).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(40400000, err, "product", primaryFilter["product_id"])
		}
		return nil, apperror.NewError(50000000, err)
	}
	return &product, nil
}

func (r *productRepository) Filter(ctx context.Context, filter *entity.Product, limit, offset int, isAsc bool) ([]entity.Product, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	query := db.WithContext(ctx).Model(&entity.Product{})
	if filter != nil {
		query = query.Where(filter)
	}

	query, appErr := r.applyListOptions(query, &entity.Product{}, limit, offset, isAsc)
	if appErr != nil {
		return nil, appErr
	}

	var products []entity.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return products, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]entity.Product, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	var products []entity.Product
	if err := db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return products, nil
}

func (r *productRepository) Delete(ctx context.Context, filter *entity.Product) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Where(filter).Delete(&entity.Product{}).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *productRepository) Count(ctx context.Context, filter *entity.Product) (int64, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	var count int64
	if err := db.WithContext(ctx).Model(&entity.Product{}).Where(filter).Count(&count).Error; err != nil {
		return 0, apperror.NewError(50000000, err)
	}
	return count, nil
}

func (r *productRepository) applyListOptions(query *gorm.DB, model any, limit, offset int, isAsc bool) (*gorm.DB, *apperror.AppError) {
	if limit < 0 {
		return nil, apperror.NewError(40000000, errors.New("limit must be greater than or equal to 0"))
	}
	if offset < 0 {
		return nil, apperror.NewError(40000000, errors.New("offset must be greater than or equal to 0"))
	}

	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(model); err != nil {
		return nil, apperror.NewError(50000000, err)
	}

	if len(stmt.Schema.PrimaryFields) > 0 {
		direction := "DESC"
		if isAsc {
			direction = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", stmt.Schema.PrimaryFields[0].DBName, direction))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return query, nil
}

func (r *productRepository) buildPrimaryFilter(filter *entity.Product) (map[string]any, *apperror.AppError) {
	if filter == nil {
		return nil, apperror.NewError(40000000, errors.New("find filter is required"))
	}

	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(&entity.Product{}); err != nil {
		return nil, apperror.NewError(50000000, err)
	}

	value := reflect.ValueOf(filter)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, apperror.NewError(40000000, errors.New("find filter is required"))
		}
		value = value.Elem()
	}

	primaryFilter := make(map[string]any, len(stmt.Schema.PrimaryFields))
	for _, primaryField := range stmt.Schema.PrimaryFields {
		fieldValue := value.FieldByName(primaryField.Name)
		if !fieldValue.IsValid() || fieldValue.IsZero() {
			return nil, apperror.NewError(40000000, fmt.Errorf("missing primary field: %s", primaryField.DBName))
		}
		primaryFilter[primaryField.DBName] = fieldValue.Interface()
	}

	return primaryFilter, nil
}
