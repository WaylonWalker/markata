{{- define "markata-notes.namespace" -}}
{{- default (printf "%s-notes" .Values.project_identifier) .Values.namespaceOverride -}}
{{- end -}}

{{- define "markata-notes.projectName" -}}
{{- default .Values.project_identifier .Values.project_name -}}
{{- end -}}

{{- define "markata-notes.labels" -}}
app.kubernetes.io/name: "markata-notes"
app.kubernetes.io/instance: "{{ .Values.project_identifier }}"
app.kubernetes.io/component: "notes"
app.kubernetes.io/part-of: "{{ include "markata-notes.projectName" . }}"
app.kubernetes.io/environment: "{{ .Values.environment | default "prod" }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
helm.sh/chart: "{{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" }}"
{{- end -}}

{{- define "markata-notes.sitePvcName" -}}
{{- default (printf "%s-notes-site%s-pvc" .Values.project_identifier (.Values.storage.site.pvcNameSuffix | default "")) .Values.storage.site.existingClaimName -}}
{{- end -}}

{{- define "markata-notes.sourcePvcName" -}}
{{- printf "%s-notes-source%s-pvc" .Values.project_identifier (.Values.storage.source.pvcNameSuffix | default "") -}}
{{- end -}}

{{- define "markata-notes.searchPvcName" -}}
{{- printf "%s-notes-search%s-pvc" .Values.project_identifier (.Values.storage.search.pvcNameSuffix | default "") -}}
{{- end -}}

{{- define "markata-notes.cachePvcName" -}}
{{- printf "%s-notes-cache%s-pvc" .Values.project_identifier (.Values.storage.cache.pvcNameSuffix | default "") -}}
{{- end -}}

{{- define "markata-notes.host" -}}
{{- default (printf "%s.example.com" .Values.project_identifier) .Values.ingress.host -}}
{{- end -}}

{{- define "markata-notes.tlsSecretName" -}}
{{- default (printf "%s-notes-tls" .Values.project_identifier) .Values.ingress.tls.secretName -}}
{{- end -}}

{{- define "markata-notes.serviceAccountName" -}}
{{- if .Values.serviceAccount.name -}}
{{- .Values.serviceAccount.name -}}
{{- else -}}
{{- printf "%s-notes-workload" .Values.project_identifier -}}
{{- end -}}
{{- end -}}

{{- define "markata-notes.validateBuilderAdminAuthInternalURL" -}}
{{- $authInternalURL := required "builderAdmin.ingress.auth.internalUrl is required when builder-admin auth is enabled" .Values.builderAdmin.ingress.auth.internalUrl -}}
{{- $isHTTPS := hasPrefix "https://" $authInternalURL -}}
{{- $isClusterLocalService := regexMatch "^http://[a-z0-9]([-a-z0-9]*[a-z0-9])?\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?\\.svc\\.cluster\\.local:(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])/?$" $authInternalURL -}}
{{- if not (or $isHTTPS $isClusterLocalService) }}{{ fail "builderAdmin.ingress.auth.internalUrl must use https:// or an http://<service>.<namespace>.svc.cluster.local:<port> URL when builder-admin auth is enabled" }}{{ end -}}
{{- end -}}
