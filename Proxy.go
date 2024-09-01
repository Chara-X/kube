package kube

import (
	core "k8s.io/api/core/v1"
)

type Proxy struct {
	*core.Service
}
