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

<a name="unreleased"></a>
### [Unreleased]
#### Feat
- **datadog:** simplify run command and integrate datadog ID adapters
- **tags:** support typed tags (string, bool, int, int64)
- **traces:** add configurable delay before closing the span

<a name="v0.0.1"></a>
### v0.0.1 - 2022-02-13
#### Feat
- **run:** add initial spike of a tracer
  
  
[Unreleased]: https://github.com/davidalpert/go-opentracer/compare/v0.0.1...HEAD
