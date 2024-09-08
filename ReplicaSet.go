package kube

import (
	"sync"

	"github.com/google/uuid"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicaSet struct {
	*apps.ReplicaSet
	pods []*Pod
}

func (r *ReplicaSet) Start(ctx *sync.Map) error {
	r.Status.Replicas = *r.Spec.Replicas
	for i := 0; i < int(*r.Spec.Replicas); i++ {
		go func() {
			for r.Status.Replicas != 0 {
				var pod = &Pod{Pod: &core.Pod{TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Pod"}, ObjectMeta: r.Spec.Template.ObjectMeta, Spec: r.Spec.Template.Spec}}
				pod.SetName(r.GetName() + "-" + uuid.New().String())
				pod.SetNamespace(r.GetNamespace())
				pod.SetCreationTimestamp(meta.Now())
				ctx.Store(pod.GetName(), pod)
				r.pods = append(r.pods, pod)
				if err := pod.Start(ctx); err == nil {
					break
				}
			}
		}()
	}
	return nil
}
func (r *ReplicaSet) Stop(ctx *sync.Map) error {
	r.Status.Replicas = 0
	for _, pod := range r.pods {
		ctx.Delete(pod.Name)
		pod.Stop(ctx)
	}
	return nil
}
