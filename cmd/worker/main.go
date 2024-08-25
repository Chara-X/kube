package main

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// func main() {
// 	var scheme = runtime.NewScheme()
// 	core.AddToScheme(scheme)
// 	var opr, ctx = operator.New(config.GetConfigOrDie(), scheme), context.Background()
// 	var inf, _ = opr.GetInformer(ctx, &core.Pod{})
// 	inf.AddEventHandler(cache.FilteringResourceEventHandler{
// 		FilterFunc: func(obj interface{}) bool { return obj.(*core.Pod).Spec.NodeName != os.Args[1] },
// 		Handler: cache.ResourceEventHandlerFuncs{
// 			AddFunc: func(obj interface{}) {
// 				exec.Command("docker", "run", "--name", obj.(*core.Pod).Spec.Containers[0].Name, obj.(*core.Pod).Spec.Containers[0].Image).Run()
// 			},
// 			UpdateFunc: func(oldObj, newObj interface{}) {
// 				var oldPod, newPod = oldObj.(*core.Pod), newObj.(*core.Pod)
// 				if oldPod.Spec.Containers[0].Image != newPod.Spec.Containers[0].Image {
// 					exec.Command("docker", "rm", newPod.Spec.Containers[0].Name).Run()
// 					exec.Command("docker", "run", "--name", newPod.Spec.Containers[0].Name, newPod.Spec.Containers[0].Image).Run()
// 				}
// 			},
// 			DeleteFunc: func(obj interface{}) {
// 				exec.Command("docker", "rm", obj.(*core.Pod).Spec.Containers[0].Name).Run()
// 			},
// 		},
// 	})
// 	opr.Start(ctx)
// }

func main() {
	var scheme = runtime.NewScheme()
	core.AddToScheme(scheme)
	var mrg, _ = manager.New(config.GetConfigOrDie(), manager.Options{
		Scheme: scheme,
	})
	controller.New("", mrg, controller.Options{})
	mrg.Start(signals.SetupSignalHandler())
}
