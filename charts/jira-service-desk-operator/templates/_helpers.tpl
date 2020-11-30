{{/*
Expand the name of the chart.
*/}}
{{- define "jira-service-desk-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "jira-service-desk-operator.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "jira-service-desk-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "jira-service-desk-operator.labels" -}}
helm.sh/chart: {{ include "jira-service-desk-operator.chart" . }}
{{ include "jira-service-desk-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "jira-service-desk-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jira-service-desk-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "jira-service-desk-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "jira-service-desk-operator.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the tls secret to use
*/}}
{{- define "jira-service-desk-operator.tlsSecretName" -}}
{{- if .Values.tlsSecretName }}
{{- .Values.tlsSecretName }}
{{- else }}
{{- printf "%s-webhook-server-cert" (include "jira-service-desk-operator.fullname" .) }}
{{- end }}
{{- end }}
