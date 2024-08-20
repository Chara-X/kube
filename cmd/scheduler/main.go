package main

import (
	"context"
	"math/rand"

	"github.com/Chara-X/kube/operator"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	var scheme = runtime.NewScheme()
	core.AddToScheme(scheme)
	var opr, ctx = operator.New(config.GetConfigOrDie(), scheme), context.Background()
	var inf, _ = opr.GetInformer(ctx, &core.Pod{})
	inf.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool { return obj.(*core.Pod).Spec.NodeName == "" },
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var pod, nodes = obj.(*core.Pod), &core.NodeList{}
				opr.List(ctx, nodes)
				pod.Spec.NodeName = nodes.Items[rand.Intn(len(nodes.Items))].Name
				opr.Update(ctx, pod)
			},
		},
	})
	opr.Start(ctx)
}
