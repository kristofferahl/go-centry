#!/usr/bin/env bash

interactive () {
  echo "interactive ($*)"
  echo -n "Enter your name: "
  read -r name
  echo "Thanks ${name:?}!"
}
