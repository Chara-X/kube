package kube

import (
	"sync"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Replicas struct {
	*apps.ReplicaSet
	pods  []*Pod
	Mutex *sync.Mutex
}

func (r *Replicas) Start(ctx *Context) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	for i := 0; i < int(*r.Spec.Replicas); i++ {
		var pod = &Pod{Pod: &core.Pod{TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Pod"}, ObjectMeta: r.Spec.Template.ObjectMeta, Spec: r.Spec.Template.Spec}}
		ctx.Create(pod.Name, pod)
		go func() {
			pod.Start(ctx)
		}()
	}
	return nil
}
func (r *Replicas) Stop(ctx *Context) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	for _, pod := range r.pods {
		ctx.Delete(pod.Name)
		go func() {
			pod.Stop(ctx)
		}()
	}
	return nil
}
