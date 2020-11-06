**Kubernetes Operator/Controller in Golang**
---
A simple k8s operator

```bash
# BUILD & RUN
$ kubectl apply -f artifacts/daemonstool/crd.yaml  # Install the CRD at Kubernetes first
#$  kubectl get crd #see all the crds installed in your k8s
#$  kubectl api-resources | grep daemonstools # check its API installed
$ go build -o ctrl .  # Build the controller (executable object)
$ ./ctrl -kubeconfig ~/.kube/config  -logtostderr=true # Run it
# Another shell, install a Daemonstool
$ kubectl apply -f artifacts/daemonstool/daemonstool-example.yaml
# Watch the pods
$ kubectl get pods -w
# Using our CRD API Resource
$ kubectl get daemonstool
```
