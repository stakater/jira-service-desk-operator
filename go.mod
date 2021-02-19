module github.com/stakater/jira-service-desk-operator

go 1.16

require (
	cloud.google.com/go v0.72.0 // indirect
	github.com/go-logr/logr v0.3.0
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/prometheus/client_golang v1.8.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/stakater/operator-utils v0.1.8
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58 // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/h2non/gock.v1 v1.0.16
	k8s.io/api v0.19.4
	k8s.io/apiextensions-apiserver v0.19.4 // indirect
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20201113171705-d219536bb9fd // indirect
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	sigs.k8s.io/controller-runtime v0.6.4
)

replace k8s.io/client-go => k8s.io/client-go v0.19.0
