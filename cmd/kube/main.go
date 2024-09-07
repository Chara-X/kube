package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/Chara-X/kube"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	var router = http.NewServeMux()
	router.HandleFunc("GET /apis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;g=apidiscovery.k8s.io;v=v2;as=APIGroupDiscoveryList")
		json.NewEncoder(w).Encode(apis)
	})
	router.HandleFunc("POST /api/v1/namespaces/default/{resource}", func(w http.ResponseWriter, r *http.Request) {
		var obj client.Object
		switch r.PathValue("resource") {
		case "pods":
			obj = &kube.Pod{}
		case "replicases":
			obj = &kube.Replicas{}
		case "ingresses":
			obj = &kube.Ingress{}
		case "configmaps":
			obj = &core.ConfigMap{}
		}
		json.NewDecoder(r.Body).Decode(obj)
		ctx.Store(obj.GetName(), obj)
		if obj, ok := obj.(kube.Object); ok {
			go func() {
				obj.Start(ctx)
			}()
		}
		json.NewEncoder(w).Encode(obj)
	})
	router.HandleFunc("DELETE /api/v1/namespaces/default/{resource}/{name}", func(w http.ResponseWriter, r *http.Request) {
		if obj, ok := ctx.LoadAndDelete(r.PathValue("name")); ok {
			if obj, ok := obj.(kube.Object); ok {
				obj.Stop(ctx)
			}
			json.NewEncoder(w).Encode(obj)
		} else {
			json.NewEncoder(w).Encode(notFound)
		}
	})
	router.HandleFunc("GET /api/v1/namespaces/default/{resource}/{name}", func(w http.ResponseWriter, r *http.Request) {
		if obj, ok := ctx.Load(r.PathValue("name")); ok {
			json.NewEncoder(w).Encode(obj)
		} else {
			json.NewEncoder(w).Encode(notFound)
		}
	})
	router.HandleFunc("GET /api/v1/namespaces/default/pods/{name}/log", func(w http.ResponseWriter, r *http.Request) {
		if obj, ok := ctx.Load(r.PathValue("name")); ok {
			for scanner := bufio.NewScanner(obj.(*kube.Pod).Stdout); scanner.Scan(); w.(http.Flusher).Flush() {
				if _, err := w.Write([]byte(scanner.Text() + "\n")); err != nil {
					break
				}
			}
		} else {
			json.NewEncoder(w).Encode(notFound)
		}
	})
	router.HandleFunc("GET /api/v1/namespaces/default/{resource}", func(w http.ResponseWriter, r *http.Request) {
		var resource, list = r.PathValue("resource"), &meta.List{TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "List"}}
		ctx.Range(func(key, value any) bool {
			var obj = value.(client.Object)
			if strings.HasPrefix(resource, strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind)) {
				list.Items = append(list.Items, runtime.RawExtension{Object: obj})
			}
			return true
		})
		json.NewEncoder(w).Encode(list)
	})
	http.ListenAndServe("127.0.0.1:6443", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req, _ = httputil.DumpRequest(r, true)
		fmt.Println(string(req))
		w.Header().Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
	}))
}
