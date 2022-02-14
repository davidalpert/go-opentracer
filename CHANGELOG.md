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


<a name="v1.1.0"></a>
### [v1.1.0] - 2022-02-13
#### Build
- read version from sbot
#### Chore
- ignore a .local scratch folder
#### Feat
- ensure that at least one valid exporter is configured
- **tags:** allow overriding the service name and version
- **tags:** allow overriding the span name

<a name="v1.0.0"></a>
### [v1.0.0] - 2022-02-13
#### Build
- add an xbuild target
- add a rebuild target
- auto-generate change logs
#### Doc
- add usage examples
#### Docs
- release notes for v1.0.0
#### Feat
- ensure that the command runs inside the outer shell environment
- **datadog:** simplify run command and integrate datadog ID adapters
- **tags:** support typed tags (string, bool, int, int64)
- **traces:** support TRACEPARENT as a formatted token
- **traces:** add configurable delay before closing the span

<a name="v0.0.1"></a>
### v0.0.1 - 2022-02-13
#### Feat
- **run:** add initial spike of a tracer

[Unreleased]: https://github.com/davidalpert/go-opentracer/compare/v1.1.0...HEAD
[v1.1.0]: https://github.com/davidalpert/go-opentracer/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/davidalpert/go-opentracer/compare/v0.0.1...v1.0.0
