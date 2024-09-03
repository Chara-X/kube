package kube

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object interface {
	client.Object
	Run(ctx context.Context) error
}
