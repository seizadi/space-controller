# space-controller

The use case is to create any resources that can not be directly created
using the application manifest. The main use case right now is to setup
secrets from vault.

The terms namespace and secret are already used by
kubernetes so I use the term "space" for this resource.

This repository implements a space controller for watching Space resources as
defined with a CustomResourceDefinition (CRD).
**Note:** There is support for only three types of secrets with the
package right now, see the configuration in the section on
vault setup for various secret types.
* Opaque
* TLS (kubernetes.io/tls)
* Certs (kubernetes.io/dockerconfigjson)

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

# list the custom resources
$ kubectl get spaces

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

To create a custom resource:
```sh
# create a CustomResourceDefinition supporting the status subresource
$ kubectl create -f artifacts/examples/crd-status-subresource.yaml
```
## Cleanup

You can clean up the created CustomResourceDefinition with:

```sh
$ kubectl create -f artifacts/examples/example-space.yaml
```
## Compatibility

HEAD of this repository will match HEAD of k8s.io/apimachinery and
k8s.io/client-go.

## Where does it come from?

`space-controller` is from
https://github.com/seizadi/space-controller.
Code changes are made in that location.

## Vault Integration

### Install Vault
If you are new to Vault here is their
[getting started](https://www.vaultproject.io/intro/getting-started/install.html)
In the following sections I assume you have Vault running. For development we run it:
```bash
$ vault server -dev
..
Unseal Key: IDjzQ/AQ4n+4UOYbD89DYziwQodetzpV3ke5+DoAQ6Y=
Root Token: 2b69b081-1a7e-7430-2027-68471114dcc6
..
```
Make sure you have VAULT_ADDR set in .profile or .bashrc
```bash
export VAULT_ADDR=http://127.0.0.1:8200
```
Then check that Vault server is running:
```bash
$ vault status
Key             Value
---             -----
Seal Type       shamir
Sealed          false
Total Shares    1
Threshold       1
Version         0.11.1
Cluster Name    vault-cluster-06f2737e
Cluster ID      0952e827-75d2-10ee-1f22-2562be2ee031
HA Enabled      false
```
Now lets setup a path for Kubernetes secrets:
```bash
$ vault secrets enable -path=k8s kv
```
You can always get list of paths or disable it:
```bash
$ vault secrets list
$ vault secrets disable k8s/
```
Now lets store a secret:
```bash
$ vault write k8s/contacts-app-seizadi-minikube-dev-secrets ATLAS_DATABASE_PASSWORD=postgres
```
We don't want to use the root-token so we create a token we can revoke,
best policy is to also attach policy to limit access to a path:
```bash
$ vault token create
token                c133f1a9-db52-145c-cd69-99ee2962f72f
...
```
You can revoke it if compromised or not needed:
```bash
$ vault token revoke cb583f98-5dce-5251-f522-3bc2012ce942
Success! Revoked token (if it existed)
```
Now we can use API request to access the secret:
```bash
$ export VAULT_TOKEN="c133f1a9-db52-145c-cd69-99ee2962f72f"
$ curl \
    --header "X-Vault-Token: $VAULT_TOKEN" \
    $VAULT_ADDR/v1/k8s/contacts-app-seizadi-minikube-dev-secrets
```
Vault setup for sample app for three types of secrets:
```bash
$ vault secrets enable -path=k8s kv
$ vault write k8s/qa0-secrets \
ATLAS_DATABASE_PASSWORD=postgres \
app-cert-tls.crt=MIIEvTCCA6WgAwIBAgIJAI1wSTI1S9DFMA0GCSqGSIb3DQEBBQUAMIGaMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExFDASBgNVBAcTC1NhbnRhIENsYXJhMREwDwYDVQQKEwhJbmZvYmxveDEMMAoGA1UECxMDQ1RPMSIwIAYDVQQDExlxYTAtdGVzdC5jc3AuaW5mb2Jsb3guY29tMSMwIQYJKoZIhvcNAQkBFhRzZWl6YWRpQGluZm9ibG94LmNvbTAeFw0xODEwMTkxNjE4NDJaFw0xOTEwMTkxNjE4NDEwJVUzELMAkGA1UECBMCQ0ExFDASBgNVBAcTC1NhbnRhIENsYXJhMREwDwYDVQQKEwhJbmZvYmxveDEMMAoGA1UECxMDQ1RPMSIwIAYDVQQDExlxYTAtdGVzdC5jc3AuaW5mb2Jsb3guY29tMSMwIQYJKoZIhvcNAQkBFhRzZWl6YWRpQGluZm9ibG94LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJkxomuER4doW53Jpt0fjm2kcQcjtCN1KCxkA7AdtE+E9BzWyFLzaL1UA7OckhvGf2ix5x7W46pTChhmNLfXRB1Yn5Rfk6CG0gt6qzGfDB3pQE3Uw63Jg7TCwTyNFqcQwU608XzVwpwOUPQVaK4Hx0fl69ulRLH2PFc5kwmIIHqun4GV1n5JUuPAkev1fkvMpPxHN716GVJ5+I4qn7vQU5Ih5N9cqV913y10xoRa0Rg+d2P8WSjJuU4/PrpovSy7RarPt6z6cAFX9SR10Ta6wHSir7HfVd7ROn4jPDrUS9AfjbyaMJzuqAHLvclEYpcsgSX2tcUO6GwcmauYd+r4a88CAwEAAaOCAQIwgf8wHQYDVR0OBBYEFKtGNXAd0oIW5IET3NgxJJFPkY/dMIHPBgNVHSMEgccwgcSAFKtGNXAd0oIW5IET3NgxJJFPkY/doYGgpIGdMIGaMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ0ExFDASBgNVBAcTC1NhbnRhIENsYXJhMREwDwYDVQQKEwhJbmZvYmxveDEMMAoGA1UECxMDQ1RPMSIwIAYDVQQDExlxYTAtdGVzdC5jc3AuaW5mb2Jsb3guY29tMSMwIQYJKoZIhvcNAQkBFhRzZWl6YWRpQGluZm9ibG94LmNvbYIJAI1wSTI1S9DFMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBABJpdkdCHyLleF6z0Ksdn6i4+6HtVM4egv1Se6HTUDnUBv2o4kEKg/S9hZIRsk0NLObFG/mCyG8q/okxpjvoK2qQlP/PHM4ElkHfUyObo1vTCRoQHakJ92PXhk1e+jMVA5hKGwVZ0x4qeO4pYNapW/Tc1I0eHD9fijCOYhDjGELV6OuIB4bEJMcYhtVexhCiMF7ozM4TZClA30oftxPSX2ij4wKvLjoJcyIDP4HdNoEHfE2+WOivgdxo1oiIIE5zkWQujQVMdNxRqYdhm+1/+JrLCke7dcCTVjLP1ZaSMRUtlcmw/kZBTawU32MiXvsR+XGhAdrvtdqulKs+IwsUqls= \
app-cert-tls.key=MIIEpAIBAAKCAQEAmTGia4RHh2hbncmm3R+ObaRxByO0I3UoLGQDsB20T4T0HNbIUvNovVQDs5ySG8Z/aLHnHtbjqlMKGGY0t9dEHViflF+ToIbSC3qrMZ8MHelATdTDrcmDtMLBPI0WpxDBTrTxfNXCnA5Q9BVorgfHR+Xr26VEsfY8VzmTCYggeq6fgZXWfklS48CR6/V+S8yk/Ec3vXoZUnn4jiqfu9BTkiHk31ypX3XfLXTGhFrRGD53Y/xZKMm5Tj8+umi9LLtFqs+3rPpwAVf1JHXRNrrAdKKvsd9V3tE6fiM8OtRL0B+NvJownO6oAcu9yURilyyBJfa1xQ7obByZq5h36vhrzwIDAQABAoIBAQCNv9ibBdYt3AlR8kIdP1LJ7yvKwGWxnXljwdOLxaCPJ+W9PZw07ReQgEnAi3LCkqRX2q2R4qLcemPP+dpz9ZMIWHWok9uE4NtAVexMSO+sSaT/n4zEpL7ipoapIZ/BTIah7lm4+g5N2g1cHOc0iOwDgiMApWbwCHkC+LouSrBK8xUVwv2LRDJCD8JjVfJJvkpD+msB9ggKF/HzpadkNYZf5MaGugYH5JvvcEB+T2GtO+ATS2utpuxUG2ov1yIArLz7i9hc8eYn+9KhxA1RuTRWNRmppjX4zQ8Wc29uvv7nNc0GGr04X9yX7lcO+IDwnYQLr/WVoivsyFLfh9nngINhAoGBAMamS4wKzx17E//T7Q2SpCH0jfidWH73GzUxpS7JVo9gmKtcZeqZjPQNWJgWSP+/B3Kk74pEPvQOOcY8p9W0Abs5Pm5+5NVuFbuMRw37Umf7FLB441lqBgSp6yWAKpyBwJl/sp6sqEJAOzU0HnOx4CI2gbp/65Ybdr9i9BJRIDcfAoGBAMVr0cWy2TEQVWKPAx5XOVJzk9kIpYKQSrKlTpPbGvxL3YTk2v1Lejz3/vfqHMbqFpqi8emmErE+c7hIQJT3UnssStdJilRmj5L6u3RT5Qq+4Assxq2Fdd0Ddd800yAKzUEP940ygKjSHzFlw1BFFFs3H0uVNrNNLk/hXNs7wKVRAoGAb0pDENYNasrFTZIBQJVi9tL3ps0gAyGVUJvbmvaZVAIeBgLh5ijYWvIPLEVv6DexiHz25lONoVVG8NSSgpsyTR2o6GaW9SuTaVsRg7fFVxPHZ4aSeEl5zasUXhILzVqz+EseWt8H9PXfNdNZLB//HavDyiRYa+Q/BsH9UzW4AqkCgYEAgSJJkLuv/bvlXhaVv57mS9x19R0GxiSD997RSz2ipS0qtObNp6lbR84f5SIpuKMeLgAvpNmQmId1QjFgrRApz4/lVHUyGosLluSTAUBvLVw1SJn9SztlITBGRb5T6z2ljM1Y6+8A4WywIquh2juVWSTxP4tWwGnXxUBwcKbhGEECgYAhDdrT4Rbzk6JdWz+YMBMtvDeHbwGGypYxNwJ0MFLKJHUN/NQ+Hqc7hAnbgDmNqoxtJOCdYRHf06A2yMXmhGX6Kyktv/xEMzKth5nMi7vGK/xfw87zfl1Y5YIXdWPC3NXR1MzetkouVoN2s2n9TRsV5iHWtc5cFwzdvh9gE056xA== \
app-image-pull.dockerconfigjson=e2F1dGhzOnt5b3VycHJpdmF0ZXJlZ2lzdHJ5LmNvbTp7dXNlcm5hbWU6amFuZWRvZSxwYXNzd29yZDp4eHh4eHh4eHh4eCxlbWFpbDpqZG9lQGV4YW1wbGUuY29tLGF1dGg6YzNSLi4uekUyfX19Cgo= \
```
***Reference ./artifacts/examples/example-space.yaml***

***Note:*** There is additional '-' for TLS added to seperate the secret
from the two keys, 'tls.crt' and 'tls.key'

***Note:*** The image pull has sensitive information I did not checkin
a valid value you will get an error if you use the above:
```bash
... error syncing 'contacts-app-seizadi-minikube-dev/app-imagepull': Secret "app-imagepull" is invalid: data[.dockerconfigjson]: Invalid value: "<secret contents redacted>": invalid character 'e' looking for beginning of value
```
[Reference the doc](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)
The value of the .dockerconfigjson field is a base64 representation of
your Docker credentials.
Create a valid secret:
```bash
kubectl create secret docker-registry regcred --docker-server=<your-registry-server> --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
```

Then use following to get the valid base64 value to store in vault:
```bash
kubectl get secret regcred --output="jsonpath={.data.\.dockerconfigjson}"
```