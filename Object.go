package kube

import "sigs.k8s.io/controller-runtime/pkg/client"

type Object interface {
	client.Object
	Start() error
	Stop() error
}
