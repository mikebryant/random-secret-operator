# Random Secret Operator
The random secret operator creates Secrets containing random data.

This is useful for frameworks like Django, that need a `SECRET_KEY` to initialise cryptography, and
needs to be the same across multiple instances.

### Install metacontroller
Follow [the upstream guide](https://metacontroller.app/guide/install/)

### Start the Operator

```bash
# Create the operator
$ kubectl apply -f deploy/code.yaml
$ kubectl apply -f deploy/operator.yaml

# Wait for the pod status to be Running
$ kubectl --namespace randomsecret get pod
NAME                                     READY     STATUS        RESTARTS   AGE
randomsecret-operator-6db7d8c7cf-l8vr4   1/1       Running       0          3s


# View the random secret CRD
$ kubectl get crd randomsecrets.randomsecrets.mikebryant.me.uk
NAME                                           AGE
randomsecrets.randomsecrets.mikebryant.me.uk   32s
```

### Create the example RandomSecret Resource
```bash
# Create the example
$ cat examples/resource.yaml
---
apiVersion: randomsecrets.mikebryant.me.uk/v1
kind: RandomSecret
metadata:
  name: minimal
  namespace: default
spec: {}

$ kubectl apply -n default -f examples/resource.yaml

# See the example resource
$ kubectl -n default get randomsecrets
NAME      AGE
minimal   12s

$ kubectl -n default get secrets minimal -o yaml
apiVersion: v1
data:
  random: NGFhZDk0YzUyNzllZmQxOTk4MjlmM2UxMzQzYjFmMGZkZTk5OTJkOWE3Mjg2YTAwMGY4MzM5YzQ4ODJjZWM1MzhhMTAyNzExODM5ZjdhMzE4OWE5YzNiNzA2NmIzYzllNzg4MDg4MTZlNjU5MDNhZDUwYTZhOGE1ODhjNzc4MTM=
kind: Secret
metadata:
  name: minimal
  namespace: default
  ownerReferences:
  - apiVersion: randomsecrets.mikebryant.me.uk/v1
    blockOwnerDeletion: true
    controller: true
    kind: RandomSecret
    name: minimal
type: Opaque

```

### Cleanup
```bash
$ kubectl delete -f examples
$ kubectl delete -f deploy
```
