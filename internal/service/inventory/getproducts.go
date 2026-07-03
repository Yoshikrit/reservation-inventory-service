package inventory

import (
	"context"

	"inventory/internal/pkg/apperror"
	"inventory/internal/pkg/json"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (s *inventoryService) GetProducts(ctx context.Context) (GetProductsResponse, *apperror.AppError) {
	var response GetProductsResponse

	data, err := s.redis.Get(ctx, cacheKeyGetProducts).Bytes()
	if err != nil && err != redis.Nil {
		log.Warn().Err(err).Msg("redis get failed, falling through to db")
	} else if err == nil {
		if errUnmarshal := json.Unmarshal(data, &response); errUnmarshal == nil {
			return response, nil
		}
	}

	productList, appErr := s.productRepo.FindAll(ctx)
	if appErr != nil {
		return response, appErr
	}

	productListResponse := make([]GetProductResponse, 0, len(productList))
	for _, productEntity := range productList {
		productListResponse = append(productListResponse, GetProductResponse{
			ProductID:   productEntity.ProductID,
			Name:        productEntity.Name,
			Description: productEntity.Description,
			Price:       productEntity.Price,
			Quantity:    productEntity.Quantity,
		})
	}

	response = GetProductsResponse{Products: productListResponse}

	if data, errMarshal := json.Marshal(response); errMarshal == nil {
		if err := s.redis.Set(ctx, cacheKeyGetProducts, data, cacheTTL).Err(); err != nil {
			log.Warn().Err(err).Msg("redis set failed")
		}
	} else {
		log.Warn().Err(errMarshal).Msg("failed to marshal products for cache")
	}

	return response, nil
}
