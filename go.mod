module github.com/stakater/jira-service-desk-operator

go 1.14

require (
	github.com/go-logr/logr v0.1.0
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/operator-framework/operator-lib v0.1.0
	github.com/operator-framework/operator-sdk v0.19.2
	github.com/stakater/operator-utils v0.1.4
	gopkg.in/h2non/gock.v1 v1.0.15
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.2
)

replace k8s.io/client-go => k8s.io/client-go v0.18.2
