module cert-manager-webhook-freenom

go 1.16

require (
	github.com/jetstack/cert-manager v1.4.0-beta.0
	github.com/tzwsoho/go-freenom v0.0.0-20201109024018-fe2c93cab446
	k8s.io/apiextensions-apiserver v0.21.1
	k8s.io/apimachinery v0.22.0-alpha.2
	k8s.io/apiserver v0.22.0-alpha.2 // indirect
	k8s.io/client-go v0.22.0-alpha.2
	k8s.io/klog/v2 v2.9.0
)
