# cert-manager-webhook-freenom
Webhook for Cert-Manager for the freenom domain provider.

The image is available on DockerHub at:
- `andreee94/cert-manager-webhook-freenom`
- `ghcr.io/andreee94/cert-manager-webhook-freenom`

## Install cert-manager

The cert-manager tested version is v1.5.0.

It is necessary to install both the `cert-manager.yaml` and the `cert-manager-crd.yaml`.

```bash
sudo kubectl apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
sudo kubectl apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager-crd.yaml
```

## Heml Chart Manifest

To generate the .yaml manifest from the helm chart contained inside the repository, 
run the following command:

```bash
make rendered-manifest.yaml
```

or install directly from github releases:

```bash
sudo kubectl apply -f https://github.com/andreee94/cert-manager-webhook-freenom/releases/latest/download/cert-manager-freenom.yaml
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

### Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: freenom
  namespace: cert-manager
type: kubernetes.io/basic-auth
stringData:
  username: REPLACE_WITH_USERNAME
  password: REPLACE_WITH_PASSWORD
```

### Cluster Issuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod-dns-issuer
spec:
  acme:
    # The ACME server URL
    server: https://acme-v02.api.letsencrypt.org/directory # https://acme-staging-v02.api.letsencrypt.org/directory
    # Email address used for ACME registration
    email: REPLACE_WITH_EMAIL
    # Name of a secret used to store the ACME account private key
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - dns01:
        webhook:
          groupName: acme.andreee94.com
          solverName: freenom
          config:
            usernameSecretRef:
              name: freenom
              key: username
            passwordSecretRef:
              name: freenom
              key: password
            ttl: 3600
            priority: 100  
```

### Wildcard Certificate

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-wildcard-certificate
  namespace: default
spec:
  dnsNames:
    - example.com
    - "*.example.com"
  secretName: example-wildcard-tls
  issuerRef:
    name: letsencrypt-prod-dns-issuer
    kind: ClusterIssuer
```

### Ingress nginx

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  namespace: default
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  tls:
    - hosts:
      - service.example.com
      secretName: example-wildcard-tls
  rules:
    - host: service.example.com
      http: 
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: example-svc
                port:
                  number: 3000
```

### Ingress traefik

In alternative to the nginx ingress it is possible to use **Treafik**.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  namespace: default
spec:
  tls:
    - hosts:
      - service.example.com
      secretName: example-wildcard-tls
  rules:
    - host: service.example.com
      http: 
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: example-svc
                port:
                  number: 3000
```


## Running the test suite

To running the test suite:

- run: `mv testdata/freenom-solver/secret.yaml.template testdata/freenom-solver/secret.yaml`
- configure the username and password inside `testdata/freenom-solver/secret.yaml`.
- run: `TEST_ZONE_NAME="example.com." make test`
- run: `make clean`


## Known Issues

While running the server is logging:

```bash
E0816 10:37:02.707993   10473 controller.go:116] loading OpenAPI spec for "v1alpha1.acme.andreee94.com" failed with: OpenAPI spec does not exist
I0816 10:37:02.708046   10473 controller.go:129] OpenAPI AggregationController: action for item v1alpha1.acme.andreee94.com: Rate Limited Requeue.
```

This warning is also present in the original repo for the creation of custom webhook challenges: [webhook-example](https://github.com/cert-manager/webhook-example/issues/27).

Please contribute in case any solution will be found.