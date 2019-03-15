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

{{- if .c.computeResources }}
resources:
{{ toYaml .c.computeResources | indent 2 }}
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