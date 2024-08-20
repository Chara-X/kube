package main

import (
	"context"
	"io"
	"math/rand"
	"net"
	"strconv"

	"github.com/Chara-X/kube/operator"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	var scheme = runtime.NewScheme()
	core.AddToScheme(scheme)
	var opr, ctx = operator.New(config.GetConfigOrDie(), scheme), context.Background()
	var inf, _ = opr.GetInformer(ctx, &core.Service{})
	inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var svc = obj.(*core.Service)
			var port = svc.Spec.Ports[0]
			var ln, _ = net.Listen("tcp", ":"+strconv.Itoa(int(port.NodePort)))
			defer ln.Close()
			for {
				var fore, _ = ln.Accept()
				var pods = &core.PodList{}
				opr.List(ctx, pods, client.InNamespace(svc.Namespace), client.MatchingLabels(svc.Spec.Selector))
				var back, _ = net.Dial("tcp", pods.Items[rand.Intn(len(pods.Items))].Status.PodIP+":"+strconv.Itoa(int(port.TargetPort.IntVal)))
				go func() {
					io.Copy(fore, back)
					back.Close()
				}()
				go func() {
					io.Copy(back, fore)
					fore.Close()
				}()
			}
		},
	})
	opr.Start(ctx)
}
