package kube

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object interface {
	client.Object
	Start(ctx *Context) error
	Stop(ctx *Context) error
}
