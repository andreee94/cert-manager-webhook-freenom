package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	// v1 "k8s.io/apiextensions-apiserver/pkg/apis/core/v1"

	//"k8s.io/client-go/kubernetes"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/jetstack/cert-manager/pkg/issuer/acme/dns/util"

	"github.com/tzwsoho/go-freenom/freenom"
)

var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	// This will register our custom DNS provider with the webhook serving
	// library, making it available as an API under the provided GroupName.
	// You can register multiple DNS provider implementations with a single
	// webhook, where the Name() method will be used to disambiguate between
	// the different implementations.
	cmd.RunWebhookServer(GroupName,
		&freenomDNSProviderSolver{},
	)
}

// freenomDNSProviderSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
// To do so, it must implement the `github.com/jetstack/cert-manager/pkg/acme/webhook.Solver`
// interface.
type freenomDNSProviderSolver struct {
	// If a Kubernetes 'clientset' is needed, you must:
	// 1. uncomment the additional `client` field in this structure below
	// 2. uncomment the "k8s.io/client-go/kubernetes" import at the top of the file
	// 3. uncomment the relevant code in the Initialize method below
	// 4. ensure your webhook's service account has the required RBAC role
	//    assigned to it for interacting with the Kubernetes APIs you need.
	client kubernetes.Clientset
}

// freenomDNSProviderConfig is a structure that is used to decode into when
// solving a DNS01 challenge.
// This information is provided by cert-manager, and may be a reference to
// additional configuration that's needed to solve the challenge for this
// particular certificate or issuer.
// This typically includes references to Secret resources containing DNS
// provider credentials, in cases where a 'multi-tenant' DNS solver is being
// created.
// If you do *not* require per-issuer or per-certificate configuration to be
// provided to your webhook, you can skip decoding altogether in favour of
// using CLI flags or similar to provide configuration.
// You should not include sensitive information here. If credentials need to
// be used by your provider here, you should reference a Kubernetes Secret
// resource and fetch these credentials using a Kubernetes clientset.
type freenomDNSProviderConfig struct {
	// Change the two fields below according to the format of the configuration
	// to be decoded.
	// These fields will be set by users in the
	// `issuer.spec.acme.dns01.providers.webhook.config` field.

	Email              string                   `json:"email"`
	PassswordSecretRef cmmeta.SecretKeySelector `json:"passwordSecretRef"`
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
// For example, `cloudflare` may be used as the name of a solver.
func (c *freenomDNSProviderSolver) Name() string {
	return "freenom"
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (c *freenomDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	fmt.Printf("Decoded configuration %v", cfg)

	password, err := c.getSecretKey(cfg.PassswordSecretRef, ch.ResourceNamespace)
	if nil != err {
		return err
	}

	err = freenom.Login(cfg.Email, password)
	if nil != err {
		return err
	}

	zone := util.UnFqdn(ch.ResolvedZone)
	fqdn := util.UnFqdn(ch.ResolvedFQDN)
	subName := fqdn[:len(fqdn)-len(zone)-1]

	err = freenom.AddRecord(ch.ResolvedZone, []freenom.DomainRecord{
		{
			Type:     freenom.RecordTypeTXT,
			Name:     subName,
			TTL:      3600,
			Value:    ch.Key, // Token to present as TXT
			Priority: 100,
		},
	})
	if nil != err {
		return err
	}

	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (c *freenomDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	// TODO: add code that deletes a record from the DNS provider's console
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	fmt.Printf("Decoded configuration %v", cfg)

	password, err := c.getSecretKey(cfg.PassswordSecretRef, ch.ResourceNamespace)
	if nil != err {
		return err
	}

	err = freenom.Login(cfg.Email, password)
	if nil != err {
		return err
	}

	zone := util.UnFqdn(ch.ResolvedZone)
	fqdn := util.UnFqdn(ch.ResolvedFQDN)
	subName := fqdn[:len(fqdn)-len(zone)-1]

	err = freenom.DeleteRecord(ch.ResolvedZone, &freenom.DomainRecord{
		Type:     freenom.RecordTypeTXT,
		Name:     subName,
		TTL:      3600,
		Value:    ch.Key, // Token to present as TXT
		Priority: 100,
	})
	if nil != err {
		return err
	}

	return nil
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (c *freenomDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	///// UNCOMMENT THE BELOW CODE TO MAKE A KUBERNETES CLIENTSET AVAILABLE TO
	///// YOUR CUSTOM DNS PROVIDER

	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	c.client = *cl

	///// END OF CODE TO MAKE KUBERNETES CLIENTSET AVAILABLE
	return nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (freenomDNSProviderConfig, error) {
	cfg := freenomDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}

// getSecretKey fetch a secret key based on a selector and a namespace
func (c *freenomDNSProviderSolver) getSecretKey(secret cmmeta.SecretKeySelector, namespace string) (string, error) {
	klog.V(6).Infof("retrieving key `%s` in secret `%s/%s`", secret.Key, namespace, secret.Name)

	sec, err := c.client.CoreV1().Secrets(namespace).Get(context.Background(), secret.Name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("secret `%s/%s` not found", namespace, secret.Name)
	}

	data, ok := sec.Data[secret.Key]
	if !ok {
		return "", fmt.Errorf("key `%q` not found in secret `%s/%s`", secret.Key, namespace, secret.Name)
	}

	return string(data), nil
}
