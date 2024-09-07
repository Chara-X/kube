package kube

import (
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object interface {
	client.Object
	Start(ctx *sync.Map) error
	Stop(ctx *sync.Map) error
}
