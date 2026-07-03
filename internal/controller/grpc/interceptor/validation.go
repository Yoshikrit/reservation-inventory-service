package interceptor

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ValidationUnary() grpc.UnaryServerInterceptor {
	v, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("protovalidate.New: %v", err))
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if msg, ok := req.(proto.Message); ok {
			if err := v.Validate(msg); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}
		return handler(ctx, req)
	}
}
