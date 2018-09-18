# space-controller

This repository implements a space controller for watching Space resources as
defined with a CustomResourceDefinition (CRD).

## References
[Extend K8 with Custom Resources](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/)

## Development
**Note:** go-get or vendor this package as `github.com/seizadi/space-controller`.

The Space controller performs operations such as:

* Register a new custom resource (custom resource type) of type `Space` using a CustomResourceDefinition.
* Create/get/list instances of your new resource type `Space`.
* Controller monitor Namesapce resource handling create/update/delete events.

It makes use of the generators in [k8s.io/code-generator](https://github.com/kubernetes/code-generator)
to generate a typed client, informers, listers and deep-copy functions. You can
do this yourself using the `./hack/update-codegen.sh` script.

The `update-codegen` script will automatically generate the following files &
directories:

* `pkg/apis/spacecontroller/v1alpha1/zz_generated.deepcopy.go`
* `pkg/client/`

Changes should not be made to these files manually, and when creating your own
controller based off of this implementation you should not copy these files and
instead run the `update-codegen` script to generate your own.

## Details

The Space controller uses [client-go library](https://github.com/kubernetes/client-go/tree/master/tools/cache) extensively.
The details of interaction points of the Space controller with various mechanisms from this library are
explained [here](docs/controller-client-go.md).


## Purpose

This Space kube-like controller creates a namespace and populates it
with secretes from a Vault.

## Running

**Prerequisite**:
   * Go version 1.10 or greater
   * Since the space-controller uses `apps/v1`
deployments, the Kubernetes cluster version should be greater than 1.9.

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
$ go build -o bin/space-controller .
$ ./bin/space-controller -kubeconfig=$HOME/.kube/config

# create a CustomResourceDefinition
$ kubectl create -f artifacts/examples/crd.yaml

# create a custom resource of type Space
$ kubectl create -f artifacts/examples/example-space.yaml

# check secrets created through the Space custom resource
$ kubectl get secrets
$ kubectl get secret example-foo -o yaml

# delete a custom resource of type Space
$ kubectl delete -f artifacts/examples/example-space.yaml

```

### Docker Build
```
$ docker build -t="soheileizadi/space-controller:latest" .
$ docker push soheileizadi/space-controller
$ kubectl create -f ./deploy/space-controller.yaml
```

## Use Cases
```
TODO >>>> Update the section below here from sample controller.
```

CustomResourceDefinitions can be used to implement custom resource types for your Kubernetes cluster.
These act like most other Resources in Kubernetes, and may be `kubectl apply`'d, etc.

Some example use cases:

* Provisioning/Management of external datastores/databases (eg. CloudSQL/RDS instances)
* Higher level abstractions around Kubernetes primitives (eg. a single Resource to define an etcd cluster, backed by a Service and a ReplicationController)

## Defining types

Each instance of your custom resource has an attached Spec, which should be defined via a `struct{}` to provide data format validation.
In practice, this Spec is arbitrary key-value data that specifies the configuration/behavior of your Resource.

For example, if you were implementing a custom resource for a Database, you might provide a DatabaseSpec like the following:

``` go
type DatabaseSpec struct {
	Databases []string `json:"databases"`
	Users     []User   `json:"users"`
	Version   string   `json:"version"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
```

## Validation

To validate custom resources, use the [`CustomResourceValidation`](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#validation) feature.

This feature is beta and enabled by default in v1.9.

### Example

The schema in [`crd-validation.yaml`](./artifacts/examples/crd-validation.yaml) applies the following validation on the custom resource:
`spec.replicas` must be an integer and must have a minimum value of 1 and a maximum value of 10.

In the above steps, use `crd-validation.yaml` to create the CRD:

```sh
# create a CustomResourceDefinition supporting validation
$ kubectl create -f artifacts/examples/crd-validation.yaml
```

## Subresources

Custom Resources support `/status` and `/scale` subresources as a [beta feature](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#subresources) in v1.11 and is enabled by default.
This feature is [alpha](https://v1-10.docs.kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#subresources) in v1.10 and to enable it you need to set the `CustomResourceSubresources` feature gate on the [kube-apiserver](https://kubernetes.io/docs/admin/kube-apiserver):

```sh
--feature-gates=CustomResourceSubresources=true
```

### Example

The CRD in [`crd-status-subresource.yaml`](./artifacts/examples/crd-status-subresource.yaml) enables the `/status` subresource
for custom resources.
This means that [`UpdateStatus`](./controller.go#L330) can be used by the controller to update only the status part of the custom resource.

To understand why only the status part of the custom resource should be updated, please refer to the [Kubernetes API conventions](https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status).

In the above steps, use `crd-status-subresource.yaml` to create the CRD:

```sh
# create a CustomResourceDefinition supporting the status subresource
$ kubectl create -f artifacts/examples/crd-status-subresource.yaml
```

## Cleanup

You can clean up the created CustomResourceDefinition with:

    $ kubectl delete crd spaces.spacecontroller.seizadi.github.com

## Compatibility

HEAD of this repository will match HEAD of k8s.io/apimachinery and
k8s.io/client-go.

## Where does it come from?

`space-controller` is from
https://github.com/seizadi/space-controller.
Code changes are made in that location.
