package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"syscall"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubernetes/pkg/apis/apidiscovery"
)

var cmds = map[string]*exec.Cmd{}

func main() {
	var router = http.NewServeMux()
	router.HandleFunc("GET /apis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList")
		var apis = &apidiscovery.APIGroupDiscoveryList{}
		var data, _ = os.ReadFile("./apis.yaml")
		yaml.Unmarshal(data, apis)
		json.NewEncoder(w).Encode(apis)
	})
	router.HandleFunc("POST /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {
		var pod = core.Pod{}
		json.NewDecoder(r.Body).Decode(&pod)
		var con = pod.Spec.Containers[0]
		var cmd = exec.Command(con.Image)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Pdeathsig: syscall.SIGKILL,
		}
		cmd.Start()
		switch pod.Spec.RestartPolicy {
		case core.RestartPolicyAlways:
			go func() {
				for {
					if err := cmd.Run(); err == nil {
						break
					}
					cmd = exec.Command(con.Image)
					cmd.Start()
				}
			}()
		case core.RestartPolicyOnFailure:
		}
		// go func() {
		// 	for state:=cmd.wa; pod.Spec.RestartPolicy==core.RestartPolicyAlways||(pod.Spec.RestartPolicy==core.RestartPolicyOnFailure&&cmd.ProcessState.ExitCode()!=0) {
		// 		cmd = exec.Command(con.Image)
		// 		cmd.Start()
		// 	}
		// }()
		cmds[con.Name] = cmd
		json.NewEncoder(w).Encode(newPod(con.Name, cmd))
	})
	router.HandleFunc("DELETE /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {
		var name = r.PathValue("name")
		var cmd = cmds[name]
		cmd.Process.Kill()
		delete(cmds, name)
		json.NewEncoder(w).Encode(newPod(name, cmd))
	})
	router.HandleFunc("GET /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {
		if cmd, ok := cmds[r.PathValue("name")]; ok {
			json.NewEncoder(w).Encode(newPod(r.PathValue("name"), cmd))
		} else {
			json.NewEncoder(w).Encode(meta.Status{Reason: "NotFound", Message: "Pod not found", Code: http.StatusNotFound})
		}
	})
	router.HandleFunc("GET /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {
		var pods = &core.PodList{TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "PodList"}}
		for k, v := range cmds {
			pods.Items = append(pods.Items, newPod(k, v))
		}
		json.NewEncoder(w).Encode(pods)
	})
	http.ListenAndServe("127.0.0.1:6443", &middleware{next: router})
}

func newPod(name string, cmd *exec.Cmd) core.Pod {
	return core.Pod{
		TypeMeta:   meta.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: meta.ObjectMeta{Name: name, Namespace: "default"},
		Spec: core.PodSpec{
			Containers: []core.Container{{Name: name, Image: cmd.Path}},
		},
		Status: core.PodStatus{Phase: core.PodPhase(cmd.ProcessState.String())},
	}
}

type middleware struct{ next http.Handler }

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dumpReq, _ = httputil.DumpRequest(r, true)
	fmt.Println(string(dumpReq))
	w.Header().Set("Content-Type", "application/json")
	m.next.ServeHTTP(w, r)
}
