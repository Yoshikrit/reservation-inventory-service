package grpcutil

import (
	"errors"

	"inventory/internal/pkg/apperror"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToAppError maps a gRPC status error to the correct AppError category.
// Add a new case here whenever apperror gains a new category that maps to a gRPC code.
func ToAppError(err error) *apperror.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return apperror.NewError(50000000, err)
	}
	msg := st.Message()
	cause := errors.New(msg)
	switch st.Code() {
	case codes.NotFound:
		return apperror.NewErrorWithDescription(apperror.CategoryNotFound, 40400000, msg, cause)
	case codes.AlreadyExists:
		return apperror.NewErrorWithDescription(apperror.CategoryConflict, 40900000, msg, cause)
	case codes.InvalidArgument, codes.FailedPrecondition:
		return apperror.NewErrorWithDescription(apperror.CategoryBadRequest, 40000000, msg, cause)
	default:
		return apperror.NewError(50000000, err)
	}
}
