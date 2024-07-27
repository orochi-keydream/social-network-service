package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const requestIdKey = "x-request-id"

func RequestIdInterceptor(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	requestId, ok := ctx.Value(requestIdKey).(string)

	if !ok {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	md, ok := metadata.FromOutgoingContext(ctx)

	if !ok {
		md = metadata.New(map[string]string{
			requestIdKey: requestId,
		})
	} else {
		md.Append(requestIdKey, requestId)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	fmt.Println(md.Get(requestIdKey)[0])

	return invoker(ctx, method, req, reply, cc, opts...)
}
