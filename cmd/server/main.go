package main

import (
	"encoding/json"
	"net/http"
)

var etcd = map[string]interface{}{}

func main() {
	http.HandleFunc("POST /api/v1/namespaces/default/pods", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("PUT /api/{group}/{version}/namespaces/{namespace}/pods/{name}", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("DELETE /api/{group}/{version}/namespaces/{namespace}/pods/{name}", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("GET /api/v1/namespaces/default/pods/{name}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(etcd[r.PathValue("name")])
	})
	http.ListenAndServe(":6443", nil)
}
