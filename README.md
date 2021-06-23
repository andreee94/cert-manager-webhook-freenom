# cert-manager-webhook-freenom
Webhook for Cert-Manager for the freenom domain provider.

The image is available on DockerHub at `andreee94/cert-manager-webhook-freenom`.

## Heml Chart Manifest

To generate the .yaml manifest from the helm chart contained inside the repository, 
run the following command:

```bash
make rendered-manifest.yaml
```

### RBAC

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: cert-manager
  name: freenom-webhook:secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  resourceNames: ["freenom"]
  verbs: ["get", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: cert-manager
  name: freenom-webhook:secret-reader
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: freenom-webhook:secret-reader
subjects:
- apiGroup: ""
  kind: ServiceAccount
  name: freenom-webhook
```


## Running the test suite

To running the test suite:

- run: `mv testdata/freenom-solver/secret.yaml.template testdata/freenom-solver/secret.yaml`
- configure the username and password inside `testdata/freenom-solver/secret.yaml`.
- run: `TEST_ZONE_NAME="example.com." make test`
- run: `make clean`
