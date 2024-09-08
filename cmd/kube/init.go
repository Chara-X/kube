package main

import (
	"os"
	"sync"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubernetes/pkg/apis/apidiscovery"
)

var notFound = meta.Status{TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Status"}, Status: "Failure", Message: "Not Found", Reason: "NotFound", Code: 404}
var ctx = &sync.Map{}
var apis = &apidiscovery.APIGroupDiscoveryList{}

func init() {
	var apisData, _ = os.ReadFile("./apis.yaml")
	yaml.Unmarshal(apisData, apis)
	var node = &core.Node{
		TypeMeta:   meta.TypeMeta{APIVersion: "v1", Kind: "Node"},
		ObjectMeta: meta.ObjectMeta{Name: "kube"},
		Spec:       core.NodeSpec{},
		Status:     core.NodeStatus{},
	}
	node.SetCreationTimestamp(meta.Now())
	ctx.Store("kube", node)
}
