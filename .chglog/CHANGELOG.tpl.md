<!-- PROJECT SHIELDS -->
<!--
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![License: GPL v3][license-shield]][license-url]
<!-- [![Issues][issues-shield]][issues-url] -->
<!-- [![Forks][forks-shield]][forks-url] -->
<!-- ![GitHub Contributors][contributors-shield] -->
<!-- ![GitHub Contributors Image][contributors-image-url] -->

<!-- PROJECT LOGO -->
<br />
<p align="center">
<h1 align="center">OpenTracer</h1>

<p align="center">
  OpenTracer is a CLI tool for wrapping shell scripts and shell commands inside an OpenTelemetry Trace and Span.
  <br />
  <a href="./README.md">README</a>
  ·
  <a href="./CHANGELOG.md"><string>CHANGELOG</string></a>
  <br />
  <a href="https://github.com/davidalpert/go-opentracer/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/go-opentracer/issues">Request Feature</a>
</p>

## Changelog

{{ if .Versions -}}
{{   if .Unreleased.CommitGroups -}}
<a name="unreleased"></a>
### [Unreleased]
{{     range .Unreleased.CommitGroups -}}
#### {{ .Title }}
{{       range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{       end -}}
{{     end -}}
{{   end -}}
{{ end -}}

{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
### {{ if .Tag.Previous }}[{{ .Tag.Name }}]{{ else }}{{ .Tag.Name }}{{ end }} - {{ datetime "2006-01-02" .Tag.Date }}
{{ range .CommitGroups -}}
#### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end -}}
{{ end -}}

{{- if .RevertCommits -}}
#### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end -}}
{{ end -}}

{{- if .MergeCommits -}}
#### Pull Requests
{{ range .MergeCommits -}}
- {{ .Header }}
{{ end -}}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
#### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{ end -}}
{{ end -}}
{{ end -}}