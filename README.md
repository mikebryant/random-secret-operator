# Random Secret Operator
The random secret operator creates Secrets containing random data.

This is useful for frameworks like Django, that need a `SECRET_KEY` to initialise cryptography, and
needs to be the same across multiple instances.

### Build
```bash
# from the root of the repo pull all the libraries needed for operator-kit (this may take a while with all the Kubernetes dependencies)
dep ensure

# build the random secret operator binary
./scripts/build

# Use the minikube environment to build
eval $(minikube docker-env)

# build the docker container
docker build -t mikebryant/random-secret-operator:0.0.1 .
```

### Start the Operator

```bash
# Create the operator
$ kubectl apply -f examples/operator.yaml

# Wait for the pod status to be Running
$ kubectl get pod -l app=random-secret-operator
NAME                                      READY     STATUS    RESTARTS   AGE
random-secret-operator-84856b9b8d-b58k7   1/1       Running   0          14m


# View the random secret CRD
$ kubectl get crd randomsecrets.randomsecrets.mikebryant.me.uk
NAME                                           AGE
randomsecrets.randomsecrets.mikebryant.me.uk   32s
```

### Create the example RandomSecret Resource
```bash
# Create the example
$ kubectl apply -f examples/resource.yaml

# See the example resource
$ kubectl get randomsecrets
NAME      AGE
minimal   12s
```

### Cleanup
```bash
kubectl delete -f examples
$ kubectl delete crd randomsecrets.randomsecrets.mikebryant.me.uk
```
