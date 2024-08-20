package main

import (
	"context"
	"crypto/rand"

	"github.com/Chara-X/kube/operator"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	var scheme = runtime.NewScheme()
	core.AddToScheme(scheme)
	var opr, ctx = operator.New(config.GetConfigOrDie(), scheme), context.Background()
	var reconcile = func(obj interface{}) {
		var rs = obj.(*apps.ReplicaSet)
		var pods = core.PodList{}
		opr.List(ctx, &pods, client.InNamespace(rs.Namespace), client.MatchingLabels(rs.Spec.Selector.MatchLabels))
		rs.Status.Replicas = int32(len(pods.Items))
		opr.Update(ctx, rs)
		for rs.Status.Replicas < *rs.Spec.Replicas {
			var name = make([]byte, 4)
			rand.Read(name)
			opr.Create(ctx, &core.Pod{ObjectMeta: meta.ObjectMeta{Namespace: rs.Namespace, Name: rs.Name + "-" + string(name), Labels: rs.Spec.Selector.MatchLabels, OwnerReferences: []meta.OwnerReference{{APIVersion: rs.APIVersion, Kind: rs.Kind, Name: rs.Name, UID: rs.UID}}}, Spec: rs.Spec.Template.Spec})
		}
		for rs.Status.Replicas > *rs.Spec.Replicas {
			opr.Delete(ctx, &core.Pod{ObjectMeta: meta.ObjectMeta{Namespace: rs.Namespace, Name: pods.Items[rs.Status.Replicas-1].Name}})
		}
	}
	var inf, _ = opr.GetInformer(ctx, &apps.ReplicaSet{})
	inf.AddEventHandler(Reconciler(reconcile))
	inf, _ = opr.GetInformer(ctx, &core.Pod{})
	inf.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool { return obj.(*core.Pod).OwnerReferences[0].Kind == "ReplicaSet" },
		Handler: Reconciler(func(obj interface{}) {
			var pod, rs = obj.(*core.Pod), &apps.ReplicaSet{}
			var ownerRef = pod.OwnerReferences[0]
			opr.Get(ctx, client.ObjectKey{Namespace: pod.Namespace, Name: ownerRef.Name}, rs)
			reconcile(rs)
		}),
	})
	opr.Start(ctx)
}

type Reconciler func(obj interface{})

func (r Reconciler) OnUpdate(oldObj, newObj interface{})         { r(newObj) }
func (r Reconciler) OnAdd(obj interface{}, isInInitialList bool) { r(obj) }
func (r Reconciler) OnDelete(obj interface{})                    { r(obj) }
