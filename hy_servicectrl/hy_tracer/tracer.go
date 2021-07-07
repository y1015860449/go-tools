package hy_tracer

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// InitTracer 初始化 jaeger Tracer
func InitTracer(srvName string, addr string) (io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: srvName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
		//config.Logger(jaeger.StdLogger),
	)

	opentracing.SetGlobalTracer(tracer) // 将jaeger tracer注册到全局
	return closer, err
}

const contextTracerKey = "Tracer-context"

func init() {
	rand.Seed(time.Now().Unix())
}

// TracerWrapper tracer 中间件
func TracerWrapper(c *gin.Context) {
	gTracer := opentracing.GlobalTracer()
	var sp opentracing.Span
	if ctx, err := gTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header)); err == nil {
		sp = gTracer.StartSpan(c.Request.URL.Path)
	} else {
		sp = gTracer.StartSpan(c.Request.URL.Path, opentracing.ChildOf(ctx))
	}
	defer sp.Finish()

	md := make(map[string]string)
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {
		log.Error(err)
	}

	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		log.Error(err)
	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	c.Set(contextTracerKey, ctx)

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	//if statusCode >= http.StatusInternalServerError {
	//	ext.Error.Set(sp, true)
	//} else if rand.Intn(100) <= nsf {
	//	ext.SamplingPriority.Set(sp, 100)
	//}

	//_ = nsf
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	}
}

// ContextWithSpan 返回context
// 为了将gin的链路信息传递下去,在gin进行rpc调用其他服务时需要转换Context
func ContextWithSpan(c *gin.Context) (ctx context.Context) {
	v, exist := c.Get(contextTracerKey)
	if !exist {
		log.Infof("ContextWithSpan not exist")
		ctx = c
		return
	}
	ctx = v.(context.Context)
	return
}
