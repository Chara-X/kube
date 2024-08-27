package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubernetes/pkg/apis/apidiscovery"
	"k8s.io/kubernetes/pkg/apis/core"
)

func main() {
	var mux = http.NewServeMux()
	mux.HandleFunc("GET /apis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList")
		var apis = &apidiscovery.APIGroupDiscoveryList{}
		var data, _ = os.ReadFile("./apis.yaml")
		yaml.Unmarshal(data, apis)
		json.NewEncoder(w).Encode(apis)
	})
	mux.HandleFunc("POST /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {
		var pod = &core.Pod{}
		json.NewDecoder(r.Body).Decode(pod)
		exec.Command(pod.Spec.Containers[0].Image).Start()
		pod.Status.Phase = core.PodRunning
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pod)
	})
	mux.HandleFunc("PUT /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {
		var pods = &core.PodList{TypeMeta: v1.TypeMeta{APIVersion: "v1", Kind: "PodList"}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pods)
	})
	http.ListenAndServe("127.0.0.1:6443", &logsMiddleware{next: mux})
}

type logsMiddleware struct{ next http.Handler }

func (l *logsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dumpReq, _ = httputil.DumpRequest(r, true)
	fmt.Println(string(dumpReq))
	l.next.ServeHTTP(w, r)
}
