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
  <a href="./README.md">README</a>
  ·
  <a href="./CHANGELOG.md">CHANGELOG</a>
  .
  <a href="./CONTRIBUTING.md"><strong>CONTRIBUTING</strong></a>
  <br />
  <a href="https://github.com/davidalpert/go-opentracer/issues">Report Bug</a>
  ·
  <a href="https://github.com/davidalpert/go-opentracer/issues">Request Feature</a>
</p>

<details open="open">
  <summary><h2 style="display: inline-block">Table of contents</h2></summary>

- [Review existing issues](#review-existing-issues)
  - [Up for grabs](#up-for-grabs)
- [Setup for local development](#setup-for-local-development)
  - [Get the Code](#get-the-code)
  - [Install prerequisites](#install-prerequisites)
    - [Using ASDF](#using-asdf)
    - [Manually](#manually)
  - [Visit the doctor](#visit-the-doctor)
  - [Run locally](#run-locally)
- [Development workflow](#development-workflow)
  - [Branch names](#branch-names)
  - [Commit message guidelines](#commit-message-guidelines)

</details>

Contributions make the open source community an great place to learn, inspire, and create.

Please review this contribution guide to streamline your experience.

## Review existing issues

Please review existing [issues](https://github.com/davidalpert/go-opentracer/issues) before reporting bug reports or requesting new features.

A quick discussion to coordinate a proposed change before you start can save hours of rework. 

### Up for grabs

The [v1.0 - feature parity](https://github.com/davidalpert/go-opentracer/projects/1) project board tracks progress through the remaining work in this current release.

## Setup for local development

### Get the Code

1. [Fork the repository on Github](https://github.com/davidalpert/go-opentracer/fork)

1. Clone your fork
   ```sh
   git clone https://github.com/your-github-name/go-opentracer.git
   ```

### Install prerequisites

#### Using ASDF

* install `asdf` - https://asdf-vm.com/guide/getting-started.html

* install required plugins

    ```
    cat .tool-versions | awk '{print $1}' | while read -r plugin_name; do asdf plugin add $plugin_name; done
    ```

* use `asdf` to install dependencies

    ```
    asdf install
    ```

#### Manually

* [Task](https://taskfile.dev/) a task runner
* [golang 1.23](https://golang.org/doc/manage-install)
  * with a working go installation:
    ```
    go install golang.org/dl/go1.23.8@latest
    go1.23.8 download
    ```
  * open a terminal with `go1.23.8` as the linked `go` binary

### Visit the doctor

This repository includes a `doctor.sh` script which validates development dependencies.

1. Verify dependencies
    ```sh
    ./.tools/doctor.sh
    ```

This script attempts to fix basic issues, for example by running `go get`.

If `doctor.sh` reports an issue that it can't resolve you may need to help it by taking action.

Please log any issues with the doctor script by [reporting a bug](https://github.com/davidalpert/go-opentracer/issues).

### Run locally

1. Build and run the tests
    ```sh
    task cit
    ```
1. Run from source
    ```sh
    go run main.go version
    go run main.go --help
    ```

## Development workflow

This project follows a standard open source fork/pull-request workflow:

1. First, [fork the repository on Github](https://github.com/davidalpert/acmob/fork)


1. Create your Feature Branch
   ```
   git checkout -b 123-amazing-feature
   ```
1. Commit your Changes
   ```
   git commit -m 'Add some AmazingFeature'
   ```
1. Make sure the code builds and all tests pass
   ```
   make cit
   ```
3. Push to the Branch
   ```
   git push origin 123-amazing-feature
   ```
4. Open a Pull Request

    https://github.com/davidalpert/acmob/compare/123-amazing-feature

### Branch names

When working on a pull request to address or resolve a Github issue, prefix the branch name with the Github issue number.

In the preceding example, after picking up an issue with an id of 123, create a branch which starts with `GH-123` or just `123-` and a hyphenated description:

```
git checkout -b 123-amazing-feature
```

### Commit message guidelines

This project uses [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) to generate [CHANGELOG](CHANGELOG.md).

Format of a conventional commit:
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

List of supported commit type tags include:
```yaml
  - "build"    # Changes that affect the build system or external dependencies
  - "ci"       # Changes to our CI configuration files and scripts 
  - "docs"     # Documentation only changes
  - "feat"     # A new feature
  - "fix"      # A bug fix
  - "perf"     # A code change that improves performance
  - "refactor" # A code change that neither fixes a bug nor adds a feature
  - "test"     # Adding missing tests or correcting existing tests
```

Prefix your commits with one of these type tags to automatically include the commit description in the [CHANGELOG](CHANGELOG.md) for the next release.

[license-shield]: https://img.shields.io/badge/License-MIT-yellow.svg
[license-url]: https://opensource.org/licenses/MIT
