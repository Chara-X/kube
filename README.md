# kube

```sh
kubectl create -f <filename> --validate=false
kubectl delete -f <filename> --wait=false
kubectl get <resource> <name> -o yaml
kubectl logs <pod-name>
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: playground
spec:
  containers:
    - image: /home/chara-x/daisy/codes/go/.experimental/playground/playground
      env:
        - name: USERNAME
          valueFrom:
            configMapKeyRef:
              name: cm
              key: username
        - name: PASSWORD
          valueFrom:
            configMapKeyRef:
              name: cm
              key: password
```

```yaml
apiVersion: v1
kind: ReplicaSet
metadata:
  name: playground
spec:
  replicas: 2
  template:
    spec:
      containers:
        - image: /home/chara-x/daisy/codes/go/.experimental/playground/playground
```

```yaml
apiVersion: v1
kind: Ingress
metadata:
  name: playground
spec:
  defaultBackend:
    service:
      port:
        number: 8080
  rules:
    - http:
        paths:
          - path: "/happy"
            backend:
              service:
                port:
                  number: 8081
          - path: "/sad"
            backend:
              service:
                port:
                  number: 8082
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm
data:
  username: "chara-x"
  password: "123"
```

## References

[Creating a Custom Scheduler in Kubernetes: A Practical Guide](https://overcast.blog/creating-a-custom-scheduler-in-kubernetes-a-practical-guide-2d9f9254f3b5?gi=b0f3b2d6b422)

[Develop on Kubernetes Series — Demystifying the For vs Owns vs Watches controller-builders in controller-runtime](https://yash-kukreja-98.medium.com/develop-on-kubernetes-series-demystifying-the-for-vs-owns-vs-watches-controller-builders-in-c11ab32a046e)

[Kube-Proxy: What is it and How it Works](https://medium.com/@amroessameldin/kube-proxy-what-is-it-and-how-it-works-6def85d9bc8f)
