package interceptor

import (
	"context"
	"errors"

	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return nil, status.Error(appErrToCode(appErr), appErr.Description)
		}

		return nil, err
	}
}

func appErrToCode(err *apperror.AppError) codes.Code {
	switch err.Category {
	case apperror.CategoryBadRequest:
		return codes.InvalidArgument
	case apperror.CategoryNotFound:
		return codes.NotFound
	case apperror.CategoryConflict:
		return codes.AlreadyExists
	case apperror.CategoryUnprocessable:
		return codes.FailedPrecondition
	case apperror.CategoryUnauthorized:
		return codes.Unauthenticated
	default:
		return codes.Internal
	}
}
