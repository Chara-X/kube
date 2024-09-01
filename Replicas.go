package kube

import (
	"sync"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Replicas struct {
	*apps.ReplicaSet
	pods []*Pod
}

func (r *Replicas) Start() error {
	var wg = sync.WaitGroup{}
	wg.Add(int(*r.Spec.Replicas))
	for i := 0; i < int(*r.Spec.Replicas); i++ {
		go func() {
			for {
				var pod = &Pod{Pod: &core.Pod{
					TypeMeta:   meta.TypeMeta{APIVersion: "v1", Kind: "Pod"},
					ObjectMeta: r.ObjectMeta,
					Spec:       r.Spec.Template.Spec,
				}}
				r.pods = append(r.pods, pod)
				if pod.Start() == nil {
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
func (r *Replicas) Stop() error {
	for _, pod := range r.pods {
		if err := pod.Stop(); err != nil {
			return err
		}
	}
	return nil
}
