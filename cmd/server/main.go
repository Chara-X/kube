package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubernetes/pkg/apis/apidiscovery"
)

var cmds = map[string]*exec.Cmd{}

func main() {
	defer func() {
		for _, cmd := range cmds {
			cmd.Process.Kill()
		}
	}()
	var mux = http.NewServeMux()
	mux.HandleFunc("GET /apis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList")
		var apis = &apidiscovery.APIGroupDiscoveryList{}
		var data, _ = os.ReadFile("./apis.yaml")
		yaml.Unmarshal(data, apis)
		json.NewEncoder(w).Encode(apis)
	})
	mux.HandleFunc("POST /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {
		var pod = core.Pod{}
		json.NewDecoder(r.Body).Decode(&pod)
		var cmd = exec.Command(pod.Spec.Containers[0].Image)
		cmd.Start()
		cmds[pod.Name] = cmd
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toPod(pod.Name, cmd))
	})
	mux.HandleFunc("PUT /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {
		var cmd, ok = cmds[r.PathValue("name")]
		w.Header().Set("Content-Type", "application/json")
		if !ok {
			var status = meta.Status{Reason: "NotFound", Message: "Pod not found", Code: http.StatusNotFound}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(status)
		} else {
			json.NewEncoder(w).Encode(toPod(r.PathValue("name"), cmd))
		}
	})
	mux.HandleFunc("GET /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {
		var pods = &core.PodList{TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "PodList"}}
		w.Header().Set("Content-Type", "application/json")
		for k, v := range cmds {
			pods.Items = append(pods.Items, toPod(k, v))
		}
		json.NewEncoder(w).Encode(pods)
	})
	http.ListenAndServe("127.0.0.1:6443", &logsMiddleware{next: mux})
}

func toPod(name string, cmd *exec.Cmd) core.Pod {
	return core.Pod{
		TypeMeta:   meta.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: meta.ObjectMeta{Name: name, Namespace: "default"},
		Spec: core.PodSpec{
			Containers: []core.Container{{Name: name, Image: cmd.Path}},
		},
		Status: core.PodStatus{Phase: core.PodPhase(cmd.ProcessState.String())},
	}
}

type logsMiddleware struct{ next http.Handler }

func (l *logsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dumpReq, _ = httputil.DumpRequest(r, true)
	fmt.Println(string(dumpReq))
	l.next.ServeHTTP(w, r)
}
