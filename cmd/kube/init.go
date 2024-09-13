package main

import (
	"os"
	"sync"

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
}
