{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "container" -}}
image: {{ .c.image.registry }}/{{ .c.image.repository }}:{{ .c.image.tag }}
imagePullPolicy: {{ .c.image.pullPolicy }}

{{- if .c.mounts }}
volumeMounts:
  {{- range $mountName, $mount := .c.mounts }}
  - name: {{ $mountName }}
    mountPath: {{ .containerPath }}
    {{- if .subPath }}
    subPath: {{ .subPath }}
    {{- end }}
  {{- end }}
{{- end }}

{{- if .c.command }}
command:
  {{- if eq "script" .c.command.type }}
{{ toYaml .c.command.interpreter | indent 2 -}}
  - |
{{ .root.Files.Get (printf "script/%s" .c.command.script) | indent 4 }}
  {{- end }}
  {{- if eq "executable" .c.command.type }}
{{ toYaml .c.command.cmd | indent 2 -}}
  {{- end }}
{{- end }}

{{- if .c.routes }}
ports:
  {{- range $portName, $route := .c.routes }}
  - name: {{ $portName }}
    containerPort: {{ .containerPort }}
    protocol: {{ .protocol }}
  {{- end }}
{{- end }}

{{- if .c.resources }}
resources:
{{ toYaml .c.resources.definition | indent 2 }}
{{- end }}

{{- if .c.probes }}
{{- if .c.probes.livenessProbe }}
livenessProbe:
{{ toYaml .c.probes.livenessProbe | indent 2 }}
{{- end }}

{{- if .c.probes.readinessProbe }}
readinessProbe:
{{ toYaml .c.probes.readinessProbe | indent 2 }}
{{- end }}
{{- end }}
{{- end -}}