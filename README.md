# cert-manager-webhook-freenom
Webhook for Cert-Manager for the freenom domain provider.

## Heml Chart Manifest

To generate the .yaml manifest from the helm chart contained inside the repository, 
run the following command:

```bash
make rendered-manifest.yaml
```


## Running the test suite

To running the test suite:

- configure the username and password inside `testdata/freenom-solver/secret.yaml.template`.
- run: `mv testdata/freenom-solver/secret.yaml.template testdata/freenom-solver/secret.yaml`
- run: `TEST_ZONE_NAME="example.com." make test`
- run: `make clean`