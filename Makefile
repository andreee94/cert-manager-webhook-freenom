OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

IMAGE_NAME := andreee94/cert-manager-webhook-freenom
IMAGE_TAG := $(shell cat .version)
CHART_VERSION := ${shell grep '^version:' deploy/freenom-webhook/Chart.yaml | egrep -o '([0-9]+.[0-9]+.[0-9]+)'}

OUT := $(shell pwd)/_out

KUBEBUILDER_VERSION=2.3.2

$(shell mkdir -p "$(OUT)")

test: checkversion _test/kubebuilder
	go mod tidy
	TEST_ASSET_ETCD="_test/kubebuilder/bin/etcd" \
	TEST_ASSET_KUBE_APISERVER="_test/kubebuilder/bin/kube-apiserver" \
	TEST_ASSET_KUBECTL="_test/kubebuilder/bin/kubectl" \
	TEST_ASSET_KUBEBUILDER="_test/kubebuilder/bin/kube-builder" \
	go test -timeout 30m -v .

_test/kubebuilder:
	curl -fsSL https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH).tar.gz -o kubebuilder-tools.tar.gz
	mkdir -p _test/kubebuilder
	tar -xvf kubebuilder-tools.tar.gz
	mv kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH)/bin _test/kubebuilder/
	rm kubebuilder-tools.tar.gz
	rm -R kubebuilder_$(KUBEBUILDER_VERSION)_$(OS)_$(ARCH)

clean: clean-kubebuilder

clean-kubebuilder:
	rm -Rf _test/kubebuilder

build: checkversion
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .
	docker build -t "$(IMAGE_NAME):latest" .
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

checkversion:
ifeq ($(CHART_VERSION), $(IMAGE_TAG))
	@echo "CHART_VERSION: $(CHART_VERSION), IMAGE_TAG: $(IMAGE_TAG) are the same"
else
	@echo "CHART_VERSION: $(CHART_VERSION), IMAGE_TAG: $(IMAGE_TAG) are different. Exiting"
	$(error CHART_VERSION: $(CHART_VERSION), IMAGE_TAG: $(IMAGE_TAG) are different. Exiting)
	exit 0
	exit 1
endif


.PHONY: rendered-manifest.yaml
rendered-manifest.yaml: checkversion
	helm template \
	    --name-template freenom-webhook \
        --set image.repository=$(IMAGE_NAME) \
        --set image.tag=v$(IMAGE_TAG) \
		--set chart.metadata.version=$(IMAGE_TAG) \
		--version=$(IMAGE_TAG) \
		--namespace cert-manager \
        deploy/freenom-webhook > "$(OUT)/rendered-manifest.yaml"
