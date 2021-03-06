#!/usr/bin/env bash

# inspired by fig's doctor.sh approach: https://github.com/withfig/config/blob/v1.0.50/tools/doctor.sh

set -e

# Output helpers
RED="\033[0;31m"
YELLOW="\033[0;33m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

pass="${GREEN}pass${NC}"
fail="${RED}fail${NC}"
BOLD=$(tput bold)
NORMAL=$(tput sgr0)

padding="...................................."

function warn {
    inlineWarn "\n$1\n"
}

function inlineWarn {
    echo -e "${YELLOW}$1${NC}"
}

function note {
    inlineNote "\n$1\n"
}

function inlineNote {
    echo -e "${CYAN}$1${NC}"
}

function command {
    echo -e "'${BLUE}$1${NC}'"
}

function contact_support {
    # echo -e "\nRun $(command "make issue") to let us know about this error!\n"
    echo -e "\nContact the maintainers to let them know about this error!\n"
}

function fix {
    # Output the command we're running to a tmp file.
    # If the command exists in the file, then we've already
    # run this command and it's likely the fix didn't
    # work and we are in an infinite loop. We should
    # exit and cleanup the file.
    if grep -q "$*" _fixes &>/dev/null; then
        rm -f _fixes
        inlineWarn "\nLooks like we've already tried this fix before and it's not working."
        contact_support
        exit
    else
        echo -e "\nmaybe we can fix this...\n"
        echo -e "> $*\n"
        echo "$*" >>_fixes
        ($*)
        # There needs to be some time for any paper util scripts to do their
        # thing. 5 seconds seems to be sufficient.
        sleep 5
        echo -e "\n${GREEN}${BOLD}fix applied.${NORMAL}${NC}"
        # Everytime we attempt a fix, there is a chance that other checks
        # will be affected. Script should be re-run to ensure we are
        # looking at an up to date environment.
        inlineNote "\nRestarting checks to see if the problem is resolved."
        (./tools/doctor.sh) && exit
    fi
}

function withPadding() {
  msg=$1
  printf "$msg %s " ${padding:${#msg}}
}

function findCmd() {
  CMD=$1
  FIX_CMD=$2
  withPadding "checking $CMD"
  set +e
  BIN=$(which "$CMD")
  BIN_EXIT_CODE=$?
  if [[ "$BIN_EXIT_CODE" == "0" ]]; then
    echo -e "found $BIN $pass"
  else
    echo -e "not found in path $fail"
    if [[ -n "$FIX_CMD" ]]; then
      fix "$FIX_CMD";
    else
      warn "no fix available"
      contact_support
      exit 1
    fi
  fi
  set -e
}

function grepVersion() {
  CMD=$1
  VERSION_CMD=$2
  PATTERN=$3
  withPadding "checking $CMD version (want $PATTERN)"
  set +e
  OUT=$($VERSION_CMD)
  BIN_EXIT_CODE=$?
  if [[ "$BIN_EXIT_CODE" == "0" && $(echo "$OUT" | grep -e "$PATTERN") ]]; then
    echo -e "$OUT $pass"
  else
    echo -e "$OUT $fail"
    printf "\n"
    exit 1
  fi
  set -e
}

# set GO version from go.mod if not already set in ENV
GO_MOD=$(cat go.mod | grep "^go" | awk '{ print $2 }')
GO_VERSION=${GO_VERSION:=$GO_MOD}

# ---------------------------------------------------------------------------------------------------------------------
#
note "inspecting dev dependencies"
findCmd make
findCmd git
findCmd go
if [[ -n "$GO_VERSION" ]]; then
  grepVersion 'go' 'go version' "$GO_VERSION"
fi

findCmd sbot "go install github.com/restechnica/semverbot/cmd/sbot@v1.0.0"
if [ ! -f .semverbot.toml ]; then
  note "initializing sbot"
  sbot init
fi

findCmd conform "go install github.com/talos-systems/conform/cmd/conform@latest"
if [ ! -f .conform.yaml ]; then
  note "initializing conform"
  cat <<-EOD > .conform.yaml
policies:
  - type: commit
    spec:
      header:
        length: 80
        imperative: true
        case: lower
        invalidLastCharacters: .
      body:
        required: false
      gpg:
        required: false
      spellcheck:
        locale: US
      maximumOfOneCommit: false
      conventional:
        types:
          - "build"    # Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
          - "ci"       # Changes to our CI configuration files and scripts (examples: CircleCi, SauceLabs)
          - "docs"     # Documentation only changes
          - "feat"     # A new feature
          - "fix"      # A bug fix
          - "perf"     # A code change that improves performance
          - "refactor" # A code change that neither fixes a bug nor adds a feature
          - "test"     # Adding missing tests or correcting existing tests
        scopes:
          - "datadog"
          - "logs"
          - "tags"
          - "traces"
        descriptionLength: 72
EOD
fi

findCmd git-chglog "go install github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.1"
if [ ! -f .chglog/config.yml ]; then
  note "initializing git-chglog"
  git-chglog --init
fi

# if we get here clean up any incomplete fixes
rm -f _fixes

# and a final newline to finish off
printf "\n"
