module github.com/stakater/jira-service-desk-operator

go 1.16

require (
	github.com/go-logr/logr v0.3.0
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/stakater/operator-utils v0.1.13
	gopkg.in/h2non/gock.v1 v1.0.16
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/kustomize/kustomize/v3 v3.8.7 // indirect
)
