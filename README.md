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
  A command-line tool to wrap shell scripts and shell commands inside an OpenTelemetry Trace and Span.
  <br />
  <a href="./README.md"><strong>README</strong></a>
  ·
  <a href="./CHANGELOG.md">CHANGELOG</a>
  .
  <a href="./CONTRIBUTING.md">CONTRIBUTING</a>
  <br />
  <a href="https://github.com/davidalpert/go-opentracer/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/go-opentracer/issues">Request Feature</a>
</p>

<details open="open">
  <summary><h2 style="display: inline-block">Table of contents</h2></summary>

- [About the project](#about-the-project)
  - [Built with](#built-with)
- [Getting started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Supported replacement tokens](#supported-replacement-tokens)
  - [Propagate traces to an OpenTelemetry-instrumented service:](#propagate-traces-to-an-opentelemetry-instrumented-service)
  - [Propagate traces to a Datadog-instrumented service:](#propagate-traces-to-a-datadog-instrumented-service)
- [Utility commands](#utility-commands)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

</details>

<!-- ABOUT THE PROJECT -->
## About the project

### Built with

* [Golang 1.17](https://golang.org/)
* [OpenTelemetry Go API and SDK](https://github.com/open-telemetry/opentelemetry-go)
* [Cobra](https://github.com/spf13/cobra)

<!-- GETTING STARTED -->
## Getting started

To get a local copy up and running follow these simple steps.

### Prerequisites

`opentracer` ships as a single-file binary so using it requires no external dependencies.

### Installation

- TBD

<!-- USAGE EXAMPLES -->
## Usage

Invoke a shell command inside an OpenTelemetry Span

```sh
opentracer -e dev --span-name RunBackup --trace-http-endpoint $OTELCOL_OTLP_HTTP_ENDPOINT /opt/backup.sh -- $(date +%F)
```

Features:
- `opentracer` performs token replacement on the command text before executing it;
- `opentracer` adds the same tokens as environment variables so any script run inside the command can also reference the trace context;
- `opentracer` automatically creates nested spans; if you use `opentracer` to run a command or script which includes another call to `opentracer` the trace context propagates through environment variables
- override the `deployment.environment` value
  - for example: `--deployment-environment dev` or `-e dev`
- add arbitrary tags with the format `--tag key:value` and opentracer adds them to the wrapping span as string values;
  - for example: `--tag client:my_company`
- add typed spans by optionally specifying one of the supported types `--tag key:value:type`
  - for example: `--tag is_registered:true:bool`
- you can send traces to any OpenTelemetry collector configured with an OTLP HTTP endpoint using `--trace-http-endpoint` or to an OpenTelemetry log file using `--trace-log-file`

### Supported replacement tokens

| Token            | Description                                                                                                                | Example                                                   |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| `TRACE_ID`       | An OpenTelemetry-formatted 128-bit hexidecimal value for the TraceID created to wrap any Spans downstream of this command. | `4bf92f3577b34da6a3ce929d0e0e4736`                        |
| `SPAN_ID`        | An OpenTelemetry-formatted 64-bit hexidecimal value for the SpanID representing the run command.                           | `00f067aa0ba902b7`                                        |
| `W3CTRACEPARENT` | The trace context for this span formatted according to the W3C [trace-context](https://w3c.github.io/trace-context/)       | `00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01` |
| `DD_TRACE_ID`    | `TRACE_ID` formatted as a 64-bit unsigned integer<br/>to conform to Datadog's `X-DATADOG-TRACE-ID` HTTP header             | `9856658736241331422`                                     |
| `DD_SPAN_ID`     | `SPAN_ID` formatted as a 64-bit unsigned integer<br/>to conform to Datadog's `X-DATADOG-PARENT-ID` HTTP header             | `1930319880373503199`                                     |

### Propagate traces to an OpenTelemetry-instrumented service:

To send the trace context downstream to an OpenTelemetry-instrumented service set the `traceparent` HTTP header which encodes the trace ID and parent span ID:
```sh
./opentracer --tag c:false -e dev --trace-http-endpoint localhost:9003 run '/usr/bin/curl -kv -H traceparent:$W3CTRACEPARENT https://your.opentelemetry-instrumented.service.com/info'
```

If you want more fine-grained control over the `traceparent` header which conforms to the W3C [trace-context](https://w3c.github.io/trace-context/) spec use the individual `TRACE_ID` and `SPAN_ID` variables:
```sh
./opentracer --tag c:false -e dev --trace-http-endpoint localhost:9003 run '/usr/bin/curl -kv -H traceparent:00-$TRACE_ID-$SPAN_ID-00 https://your.opentelemetry-instrumented.service.com/info'
```

### Propagate traces to a Datadog-instrumented service:

Datadog uses a proprietary format for trace and parent IDs. If you want to propagate trace context to a datadog-instrumented service appropriately formatted DD_TRACE_ID and DD_SPAN_ID tokens also available:
```sh
./opentracer --tag c:134:int -e dev --trace-http-endpoint localhost:9003 run '/usr/bin/curl -kv -H X-DATADOG-TRACE-ID:$DD_TRACE_ID -H X-DATADOG-PARENT-ID:$DD_SPAN_ID https://your.datadog-instrumented.service.com/info'
```

## Utility commands

The `opentracer` binary also ships with utility commands which you can explore using the `--help` flag:

```
$> bin/opentracer
opentracer executes a shell command in an open trace

Usage:
  opentracer [command]

Available Commands:
  help        Help about any command
  run         runs a command inside an open trace and span
  version     Show version information

Flags:
  -h, --help   help for opentracer

Use "opentracer [command] --help" for more information about a command.
```

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/davidalpert/go-opentracer/issues) for a list of proposed features and known issues.

<!-- CONTRIBUTING -->
## Contributing

See the [contributor guide](https://github.com/davidalpert/go-opentracer/blob/main/CONTRIBUTING.md)

<!-- LICENSE -->
## License

Distributed under the GPU v3 License. See [LICENSE](LICENSE) for more information.

<!-- CONTACT -->
## Contact

Project Link: [https://github.com/davidalpert/go-opentracer](https://github.com/davidalpert/go-opentracer)

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/davidalpert/go-opentracer
[contributors-image-url]: https://contrib.rocks/image?repo=davidalpert/go-opentracer
[forks-shield]: https://img.shields.io/github/forks/davidalpert/go-opentracer
[forks-url]: https://github.com/davidalpert/go-opentracer/network/members
[issues-shield]: https://img.shields.io/github/issues/davidalpert/go-opentracer
[issues-url]: https://github.com/davidalpert/go-opentracego-opentracer
[license-shield]: https://img.shields.io/badge/License-GPLv3-blue.svg
[license-url]: https://www.gnu.org/licenses/gpl-3.0
