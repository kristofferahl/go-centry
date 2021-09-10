#!/usr/bin/env bash

optiontest:args() {
  echo "args ($*)"
}

optiontest:printenv() {
  env | sort
}

optiontest:noop() {
  return 0
}

# centry.cmd[optiontest:required].option[abc]/required=true
# centry.cmd[optiontest:required].option[def]/required=false
optiontest:required() {
  echo "This command should not run without required options specified..."
  env | sort
}
