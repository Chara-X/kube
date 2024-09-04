package kube

import (
	"context"
	"net/http"
	"net/http/httputil"
	"strconv"

	networking "k8s.io/api/networking/v1"
)

type Ingress struct {
	*networking.Ingress
	srv *http.Server
}

func (i *Ingress) Start(ctx Context) error {
	i.srv = &http.Server{
		Addr: ":" + strconv.Itoa(int(i.Spec.DefaultBackend.Service.Port.Number)),
		Handler: &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				for _, path := range i.Spec.Rules[0].HTTP.Paths {
					if req.URL.Path == path.Path {
						req.URL.Host = req.URL.Hostname() + ":" + strconv.Itoa(int(path.Backend.Service.Port.Number))
					}
				}
			},
		},
	}
	return i.srv.ListenAndServe()
}

func (i *Ingress) Stop(ctx *Context) error {
	return i.srv.Shutdown(context.Background())
}
