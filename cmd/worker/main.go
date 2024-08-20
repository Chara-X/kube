package main

import (
	"context"
	"os"
	"os/exec"

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
		FilterFunc: func(obj interface{}) bool { return obj.(*core.Pod).Spec.NodeName != os.Args[1] },
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				exec.Command("docker", "run", "--name", obj.(*core.Pod).Spec.Containers[0].Name, obj.(*core.Pod).Spec.Containers[0].Image).Run()
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				var oldPod, newPod = oldObj.(*core.Pod), newObj.(*core.Pod)
				if oldPod.Spec.Containers[0].Image != newPod.Spec.Containers[0].Image {
					exec.Command("docker", "rm", newPod.Spec.Containers[0].Name).Run()
					exec.Command("docker", "run", "--name", newPod.Spec.Containers[0].Name, newPod.Spec.Containers[0].Image).Run()
				}
			},
			DeleteFunc: func(obj interface{}) {
				exec.Command("docker", "rm", obj.(*core.Pod).Spec.Containers[0].Name).Run()
			},
		},
	})
	opr.Start(ctx)
}
