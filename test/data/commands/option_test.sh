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
