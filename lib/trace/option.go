package trace

import (
	"net/url"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Option interface {
	apply(opt *option)
}

type optionFunc func(opt *option)

func (fn optionFunc) apply(opt *option) {
	fn(opt)
}

type option struct {
	enabled         bool
	agentAddr       string
	serviceName     string
	environment     string
	version         string
	sdktraceOptions []sdktrace.TracerProviderOption
}

var defaultOption = &option{
	serviceName: "github.com/moneyforward/apex-common-internal-token",
	environment: "local",
}

func WithEnabled(enabled bool) Option {
	return optionFunc(func(opt *option) {
		opt.enabled = enabled
	})
}

func WithAgentAddr(addr string) Option {
	return optionFunc(func(opt *option) {
		opt.agentAddr = addr
	})
}

func WithServiceName(name string) Option {
	return optionFunc(func(opt *option) {
		opt.serviceName = name
	})
}

func WithEnvironment(env string) Option {
	return optionFunc(func(opt *option) {
		opt.environment = env
	})
}

func WithVersion(version string) Option {
	return optionFunc(func(opt *option) {
		opt.version = version
	})
}

func WithSdkTraceOption(opts ...sdktrace.TracerProviderOption) Option {
	return optionFunc(func(opt *option) {
		opt.sdktraceOptions = append(opt.sdktraceOptions, opts...)
	})
}

type urlJoinPathFunc func(base string, elem ...string) (result string, err error)

var (
	_ urlJoinPathFunc = url.JoinPath
)

type startSpanOption struct {
	attributes      []attribute.KeyValue
	urlJoinPathFunc urlJoinPathFunc
}

type StartSpanOption interface {
	apply(opt *startSpanOption)
}

type startSpanOptionFunc func(opt *startSpanOption)

func (fn startSpanOptionFunc) apply(opt *startSpanOption) {
	fn(opt)
}

func WithAttribute(attrs ...attribute.KeyValue) StartSpanOption {
	return startSpanOptionFunc(func(opt *startSpanOption) {
		opt.attributes = attrs
	})
}

func WithUrlJoinPathFunc(fn urlJoinPathFunc) StartSpanOption {
	return startSpanOptionFunc(func(opt *startSpanOption) {
		opt.urlJoinPathFunc = fn
	})
}
