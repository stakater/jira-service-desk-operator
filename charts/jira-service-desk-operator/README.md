# jira-service-desk-operator

A Helm chart to deploy jira-service-desk-operator

## Pre-requisites

- Make sure that [certman](https://cert-manager.io/) is deployed in your cluster since webhooks require certman to generate valid certs since webhooks serve using HTTPS

```terminal
$ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.1/cert-manager.yaml
```

## Installing the chart

Helm doesn't support templatization and upgrade or deletion for CRDs. We mantain them in a separate chart which needs to be installed before you install the operator.

```sh
helm repo add stakater https://stakater.github.io/stakater-charts/
helm repo update
helm install stakater/jira-service-desk-operator --namespace jira-service-desk-operator
```