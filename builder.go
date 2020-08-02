package bjaeger

import (
	"container/list"

	"github.com/go-masonry/mortar/interfaces/trace"
	"github.com/uber/jaeger-client-go/config"
)

type JaegerBuilder interface {
	SetServiceName(name string) JaegerBuilder
	AddOptions(options ...config.Option) JaegerBuilder
	SetCustomConfig(custom *config.Configuration) JaegerBuilder
	Build() (trace.OpenTracer, error)
}

type jaegerConfig struct {
	name    string
	options []config.Option
	conf    *config.Configuration
}

type jaegerBuilder struct {
	ll *list.List
}

func Builder() JaegerBuilder {
	return &jaegerBuilder{
		ll: list.New(),
	}
}

func (jb *jaegerBuilder) SetServiceName(name string) JaegerBuilder {
	jb.ll.PushBack(func(cfg *jaegerConfig) {
		cfg.name = name
	})
	return jb
}

func (jb *jaegerBuilder) AddOptions(options ...config.Option) JaegerBuilder {
	jb.ll.PushBack(func(cfg *jaegerConfig) {
		cfg.options = append(cfg.options, options...)
	})
	return jb
}
func (jb *jaegerBuilder) SetCustomConfig(custom *config.Configuration) JaegerBuilder {
	jb.ll.PushBack(func(cfg *jaegerConfig) {
		cfg.conf = custom
	})
	return jb
}

func (jb *jaegerBuilder) Build() (instance trace.OpenTracer, err error) {
	cfg := new(jaegerConfig)
	for e := jb.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(config *jaegerConfig))
		f(cfg)
	}
	if cfg.conf == nil {
		if cfg.conf, err = config.FromEnv(); err != nil {
			return nil, err
		}
	}
	if len(cfg.name) > 0 {
		cfg.conf.ServiceName = cfg.name // set name
	}
	return newWrapper(cfg), nil
}
