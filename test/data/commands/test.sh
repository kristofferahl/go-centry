#!/usr/bin/env bash

test:args() {
  echo "test:args ($*)"
}

test:env() {
  env | sort
}
