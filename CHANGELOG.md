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


<a name="v2.0.0"></a>
### [v2.0.0] - 2022-02-19
#### Build
- run tests before updating changelog and tagging a release
- refactor the ship target
- add a 'doctor' target
#### Chore
- hide completion command
#### Docs
- add an explicit run step for local development
- remove outdated acknowledgements (from another project)
#### Feat
- propagate opentrace trace context to nested invocations
- add a --debug flag which will dump some diagnostics to the console
#### Fix
- improve predictability of command parsing
#### Refactor
- rename gopentracer to opentracer
#### BREAKING CHANGE

existing run args should be reviewed and updated for compatibility with improved arg parsing

<a name="v1.1.1"></a>
### [v1.1.1] - 2022-02-14
#### Build
- ensure that the version file is up-to-date during a cross build
#### Docs
- release notes for v1.1.1
#### Fix
- set error on span when run command has non-zero exit code

<a name="v1.1.0"></a>
### [v1.1.0] - 2022-02-13
#### Build
- read version from sbot
#### Chore
- ignore a .local scratch folder
#### Docs
- release notes for v1.1.0
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

[Unreleased]: https://github.com/davidalpert/go-opentracer/compare/v2.0.0...HEAD
[v2.0.0]: https://github.com/davidalpert/go-opentracer/compare/v1.1.1...v2.0.0
[v1.1.1]: https://github.com/davidalpert/go-opentracer/compare/v1.1.0...v1.1.1
[v1.1.0]: https://github.com/davidalpert/go-opentracer/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/davidalpert/go-opentracer/compare/v0.0.1...v1.0.0
