## code

```go
paramsToRemove := []string{
{{- range .ParamsToRemove}}
   "{{.}}",
{{- end}}
}
```


# Remaining Parameters Report

{{range .ParamInfoList}}

## Parameter: {{.Param}}
Source URL:

{{.SourceURL}}

File:

{{.FilePath}}

[Search Google]({{.SearchURL}})

{{end}}

