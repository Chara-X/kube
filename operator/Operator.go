package operator

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Operator struct {
	client.Reader
	client.Writer
	cache.Informers
}

func New(config *rest.Config, scheme *runtime.Scheme) *Operator {
	var cache, _ = cache.New(config, cache.Options{Scheme: scheme})
	var client, _ = client.New(config, client.Options{Scheme: scheme, Cache: &client.CacheOptions{
		Reader: cache,
	}})
	return &Operator{cache, client, cache}
}
