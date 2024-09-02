package kube

import (
	core "k8s.io/api/core/v1"
)

type ConfigMap struct {
	*core.ConfigMap
}

func (cm *ConfigMap) Start() error { return nil }
func (cm *ConfigMap) Stop() error  { return nil }
