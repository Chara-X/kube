package kube

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"

	networking "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Ingress struct {
	*networking.Ingress
	srv *http.Server
}

func (i *Ingress) Start(ctx *sync.Map) error {
	i.SetCreationTimestamp(meta.Now())
	i.srv = &http.Server{
		Addr: ":" + strconv.Itoa(int(i.Spec.DefaultBackend.Service.Port.Number)),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range i.Spec.Rules[0].HTTP.Paths {
				if r.URL.Path == path.Path {
					httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "127.0.0.1:" + strconv.Itoa(int(path.Backend.Service.Port.Number)), Path: r.URL.Path, RawQuery: r.URL.RawQuery}).ServeHTTP(w, r)
					return
				}
			}
			http.NotFound(w, r)
		}),
	}
	return i.srv.ListenAndServe()
}

func (i *Ingress) Stop(ctx *sync.Map) error {
	return i.srv.Shutdown(context.Background())
}
