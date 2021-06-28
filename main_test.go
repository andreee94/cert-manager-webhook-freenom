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
		dns.SetBinariesPath("_test/kubebuilder/bin"),
		dns.SetUseAuthoritative(false),
		dns.SetResolvedFQDN(fmt.Sprintf("cert-manager-dns%d-tests.%s", randomIndex, zone)), // randomize FQDN to avoid frequent update problems
		dns.SetPropagationLimit(10*time.Minute),
	)

	// solver := example.New("59351")
	// fixture := dns.NewFixture(solver,
	// 	dns.SetResolvedZone("example.com."),
	// 	dns.SetManifestPath("testdata/my-custom-solver"),
	// 	dns.SetBinariesPath("_test/kubebuilder/bin"),
	// 	dns.SetDNSServer("127.0.0.1:59351"),
	// 	dns.SetUseAuthoritative(false),
	// )

	// kubectlgetnamespaces()

	fixture.RunConformance(t)
}

// var DefaultKubeAPIServerFlags = []string{
// 	"--etcd-servers={{ if .EtcdURL }}{{ .EtcdURL.String }}{{ end }}",
// 	"--cert-dir={{ .CertDir }}",
// 	"--insecure-port={{ if .URL }}{{ .URL.Port }}{{ end }}",
// 	"--insecure-bind-address={{ if .URL }}{{ .URL.Hostname }}{{ end }}",
// 	"--secure-port={{ if .SecurePort }}{{ .SecurePort }}{{ end }}",
// 	"--admission-control=AlwaysAdmit",
// }

// func kubectlgetnamespaces() {
// 	binariesPath := "_test/kubebuilder/bin"
// 	kubectlManifestsPath := "testdata/freenom-solver"
// 	name := "test"

// 	controlPlane := &integration.ControlPlane{}
// 	controlPlane.APIServer = &integration.APIServer{
// 		Args: DefaultKubeAPIServerFlags,
// 		Path: binariesPath + "/kube-apiserver",
// 	}
// 	controlPlane.Etcd = &integration.Etcd{
// 		Path: binariesPath + "/etcd",
// 	}
// 	if err := controlPlane.Start(); err != nil {
// 		fmt.Printf("error starting apiserver: %v", err)
// 	}
// 	fmt.Printf("started apiserver on %q", controlPlane.APIURL())
// 	// Create the *rest.Config for creating new clients
// 	restConfig := &rest.Config{
// 		Host: controlPlane.APIURL().Host,
// 		// gotta go fast during tests -- we don't really care about overwhelming our test API server
// 		QPS:   1000.0,
// 		Burst: 2000.0,
// 	}
// 	// var err error
// 	if clientset, err := kubernetes.NewForConfig(restConfig); err != nil {
// 		_ = clientset
// 		fmt.Printf("error constructing clientset: %v", err)
// 	}

// 	kubectl := controlPlane.KubeCtl()
// 	kubectl.Path = binariesPath + "/kubectl"

// 	// stdout, stderr, err := kubectl.Run("describe", "namespaces")

// 	// buf := new(strings.Builder)
// 	// _, _ = io.Copy(buf, stdout)
// 	// fmt.Println("stdout: ", buf.String())

// 	// buf2 := new(strings.Builder)
// 	// _, _ = io.Copy(buf2, stderr)
// 	// fmt.Println("stdout: ", buf2.String())

// 	// log.Printf("kubectl describe namespaces err: %v", err)

// 	if err := filepath.Walk(kubectlManifestsPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() || filepath.Base(path) == "config.json" {
// 			return nil
// 		}

// 		switch filepath.Ext(path) {
// 		case ".json", ".yaml", ".yml":
// 		default:
// 			fmt.Printf("skipping file %q with unrecognised extension", path)
// 			return nil
// 		}

// 		content, _ := ioutil.ReadFile(path)
// 		fmt.Printf("\ncontent of %v: %v\n", path, string(content))

// 		stdout, stderr, err := kubectl.Run("apply", "--namespace", name, "-f", path)

// 		buf := new(strings.Builder)
// 		_, _ = io.Copy(buf, stdout)
// 		fmt.Println("stdout: ", buf.String())

// 		buf2 := new(strings.Builder)
// 		_, _ = io.Copy(buf2, stderr)
// 		fmt.Println("stdout: ", buf2.String())

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Printf("created fixture %q", name)
// 		return nil
// 	}); err != nil {
// 		fmt.Printf("error creating test fixtures: %v", err)
// 	}

// 	// stopCh := make(chan struct{})
// 	// testSolver.Initialize(restConfig, stopCh)
// }
