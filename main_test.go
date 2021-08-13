package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/jetstack/cert-manager/test/acme/dns"
	// "github.com/cert-manager/webhook-example/example"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(1000)

	// Uncomment the below fixture when implementing your custom DNS provider
	fixture := dns.NewFixture(&freenomDNSProviderSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/freenom-solver"),
		// dns.SetBinariesPath("_test/kubebuilder/bin"),
		dns.SetPollInterval(10*time.Second),
		dns.SetUseAuthoritative(true),
		dns.SetResolvedFQDN(fmt.Sprintf("cert-manager-dns%d-tests.%s", randomIndex, zone)), // randomize FQDN to avoid frequent update problems
		dns.SetPropagationLimit(15*time.Minute),
	)

	fixture.RunConformance(t)
}
