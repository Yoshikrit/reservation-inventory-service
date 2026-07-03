package inventory

import (
	"context"

	"inventory/internal/entity"
	"inventory/internal/pkg/apperror"
	"inventory/internal/pkg/json"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (s *inventoryService) GetProductByID(ctx context.Context, request *GetProductRequest) (GetProductResponse, *apperror.AppError) {
	var response GetProductResponse
	cacheKey := cacheKeyProductByID(request.ProductID)

	data, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err != nil && err != redis.Nil {
		log.Warn().Err(err).Msg("redis get failed, falling through to db")
	} else if err == nil {
		if errUnmarshal := json.Unmarshal(data, &response); errUnmarshal == nil {
			return response, nil
		}
	}

	product, appErr := s.productRepo.Find(ctx, &entity.Product{ProductID: request.ProductID})
	if appErr != nil {
		return response, appErr
	}

	response = GetProductResponse{
		ProductID:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
	}

	if data, errMarshal := json.Marshal(response); errMarshal == nil {
		if err := s.redis.Set(ctx, cacheKey, data, cacheTTL).Err(); err != nil {
			log.Warn().Err(err).Msg("redis set failed")
		}
	} else {
		log.Warn().Err(errMarshal).Msg("failed to marshal product for cache")
	}

	return response, nil
}
