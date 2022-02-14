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
  <a href="./README.md"><strong>README</strong></a>
  ·
  <a href="./CHANGELOG.md">CHANGELOG</a>
  <br />
  <a href="https://github.com/davidalpert/go-opentracer/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/go-opentracer/issues">Request Feature</a>
</p>

<details open="open">
  <summary><h2 style="display: inline-block">Table of Contents</h2></summary>

- [About The Project](#about-the-project)
    - [Built With](#built-with)
- [Getting Started](#getting-started)
    - [Installation](#installation)
- [Usage](#usage)
    - [Utility Commands](#utility-commands)
- [Roadmap](#roadmap)
- [Local Development](#local-development)
    - [Prerequisites](#prerequisites)
    - [Make targets](#make-targets)
    - [Rapid Development](#rapid-development)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgements](#acknowledgements)

</details>

<!-- ABOUT THE PROJECT -->
## About The Project

### Built With

* [Golang 1.17](https://golang.org/)
* [OpenTelemetry Go API and SDK](https://github.com/open-telemetry/opentelemetry-go)
* [Cobra](https://github.com/spf13/cobra)

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

`opentracer` is built and distributed as a single-file binary so there are no prerequisites.

### Installation

- TBD

<!-- USAGE EXAMPLES -->
## Usage

Run the `opentracer` binary with no arguments to get help text.

### Utility Commands

The `opentracer` binary also ships with a number of utility commands which you can explore using the `--help` flag:

```
$> bin/opentracer
gopentracer executes a shell command in an open trace

Usage:
  gopentracer [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  run         runs a command inside an open trace and span
  version     Show version information

Flags:
  -h, --help   help for gopentracer

Use "gopentracer [command] --help" for more information about a command.

```

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/davidalpert/go-opentracer/issues) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->
## Local Development

### Prerequisites

* [golang](https://golang.org/doc/manage-install)
    * with a working go installation:
      ```
      go install golang.org/dl/go1.17@latest
      go1.17 download
      ```
* [make](https://www.gnu.org/software/make/manual/html_node/index.html#Top) (often comes pre-installed or installed with other dev tooling)

Then:

1. Clone the repo
   ```sh
   git clone https://github.com/davidalpert/go-opentracer.git
   ```

2. Run the setup script
    ```sh
    ./tools/setup.sh
    ```

3. Run tests
    ```sh
    make test
    ```
   
4. Build the tool
    ```sh
    make build
    ```

### Make targets

This repo includes a `Makefile` for help running common tasks.

Run `make` with no args to list the available targets:
```
 ❯ make

  0.0.1 - available targets:

build                          build
changelog                      Generate/update CHANGELOG.md
gen                            invoke go generate
rebuild                        rebuild
run                            run direct from source
test-verbose                   run all tests (with verbose flag)
test                           run all tests
tidy                           runs 'go mod tidy' with the current versioned go command
----------                     ------------------
release-major                  release major version
release-minor                  release minor version
release-patch                  release patch version

```

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->
## License

Distributed under the GPU v3 License. See [LICENSE](LICENSE) for more information.

<!-- CONTACT -->
## Contact

David Alpert - [@davidalpert](https://twitter.com/davidalpert)

Project Link: [https://github.com/davidalpert/go-opentracer](https://github.com/davidalpert/go-opentracer)

<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements

* [David Fowler](https://github.com/davidfowl) and [Phillip Seith](https://github.com/philippseith) for the [GoLang implementation](https://github.com/philippseith/signalr) of the server-side SignalR protocol
* [Aaron Lindsay](https://github.com/aclindsa) for [ofxgo](https://github.com/aclindsa/ofxgo) (an OFX parsing library)

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