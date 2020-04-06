#!/usr/bin/env bash

get() {
  echo "get ($*)"
}

# centry.cmd[get:sub]/description=Description for subcommand
# centry.cmd[get:sub]/help=Help text for sub command
get:sub() {
  echo "get:sub ($*)"
}
