{{- define "configs" -}}
{{- $context := . }}
{{- range $path, $bytes := .Files.Glob "config/*" }}
{{ printf "%s" $path | replace "." "_" | replace "config/" "" }}: "{{ tpl $context.Values.configPath $context }}/{{ $path | replace "config/" "" }}"
{{- end }}
{{- end -}}

{{- define "secrets" -}}
{{- $context := . -}}
{{- range $path, $bytes := .Files.Glob "secrets/*" }}
{{ printf "%s" $path | replace "." "_" | replace "secrets/" "" }}: "{{ tpl $context.Values.secretsPath $context }}/{{ $path | replace "secrets/" "" }}"
{{- end }}
{{- end -}}

{{- define "appconfig" -}}
{{- $name := .Chart.Name -}}

app.yaml: |
  ports:
  {{- range .Values.containers }}
    {{ default .name $name }}:
      {{- range .routes }}
      {{ .name }}: "{{ .containerPort }}"
      {{- end }}
  {{- end }}
  dependencies:
{{ toYaml .Values.dependencies | indent 4 }}
  configs:
{{ include "configs" . | default "{}" | indent 4 }}
  secrets:
{{ include "secrets" . | default "{}" | indent 4 }}
{{- end -}}

{{- define "appcredentials" -}}
credentials.yaml: {{ toYaml .Values.credentials | b64enc }}
{{- end -}}

{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

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

{{- define "chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}