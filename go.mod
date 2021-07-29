module github.com/darkowlzz/crd-api-conversion-demo

go 1.16

require (
	// github.com/darkowlzz/operator-toolkit v0.0.0-20210727040014-66d8dab122b3
	github.com/darkowlzz/operator-toolkit v0.0.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	k8s.io/api v0.20.2
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/darkowlzz/operator-toolkit => ../operator-toolkit/
