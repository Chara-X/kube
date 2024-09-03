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

func (ig *Ingress) Run(ctx context.Context) error {
	ig.srv = &http.Server{
		Addr: ":" + strconv.Itoa(int(ig.Spec.DefaultBackend.Service.Port.Number)),
		Handler: &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				for _, path := range ig.Spec.Rules[0].HTTP.Paths {
					if req.URL.Path == path.Path {
						req.URL.Host = req.URL.Hostname() + ":" + strconv.Itoa(int(path.Backend.Service.Port.Number))
					}
				}
			},
		},
	}
	go func() {
		<-ctx.Done()
		ig.srv.Shutdown(context.Background())
	}()
	return ig.srv.ListenAndServe()
}
