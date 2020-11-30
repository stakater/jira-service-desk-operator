# managed-openshift-operator

A Helm chart to deploy jira-service-desk-operator

## Installing the chart

```sh
helm repo add stakater https://stakater.github.io/stakater-charts/
helm repo update
helm install stakater/jira-service-desk-operator --namespace jira-service-desk-operator
```

## Known Issues

- Helm doesn't support upgrading or deleting CRDs. That needs to be done manually.

```sh
kubectl apply -f charts/jira-service-desk-operator/crds
```
