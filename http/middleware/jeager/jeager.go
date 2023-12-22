package jeager

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func WithTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var span opentracing.Span
		var tracer = opentracing.GlobalTracer()

		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			// 如果上游没有span，这里重新开启一个span
			span = tracer.StartSpan(c.Request.URL.Path)
		} else {
			// 如果上游已经有span了，这里继承自上游span
			span = opentracing.StartSpan(c.Request.URL.Path, opentracing.ChildOf(spCtx))
		}
		if span != nil {
			defer func() {
				span.Finish()
			}()
		}

		c.Next()
	}
}

type MD struct {
	m metadata.MD
}

func (m *MD) ForeachKey(handler func(key, val string) error) error {
	if m.m.Len() > 0 {
		for k, vs := range m.m {
			for _, v := range vs {
				if err := handler(k, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *MD) Set(key string, value string) {
	if m.m == nil {
		m.m = metadata.New(map[string]string{})
	}
	m.m.Set(key, value)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var tracer = opentracing.GlobalTracer()
		var spanCtx context.Context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		eCxt, err := tracer.Extract(opentracing.TextMap, &MD{m: md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			return nil, err
		} else {
			span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(eCxt), ext.SpanKindRPCServer)
			defer span.Finish()
			spanCtx = opentracing.ContextWithSpan(ctx, span)
		}
		return handler(spanCtx, req)
	}
}

func UnaryClientInterceptor(c context.Context) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var tracer = opentracing.GlobalTracer()
		var span, _ = opentracing.StartSpanFromContext(c, method, ext.SpanKindRPCClient)

		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		err := tracer.Inject(span.Context(), opentracing.TextMap, &MD{m: md})
		if err != nil {
			span.LogFields(log.String("inject-err", err.Error()))
		}

		newCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.LogFields(log.String("call-error", err.Error()))
		}
		return err
	}
}
