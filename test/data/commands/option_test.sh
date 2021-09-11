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

# centry.cmd[optiontest:required].option[stringopt]/required=true
# centry.cmd[optiontest:required].option[boolopt]/type=bool
# centry.cmd[optiontest:required].option[boolopt]/required=true
# centry.cmd[optiontest:required].option[selectopt1]/required=true
# centry.cmd[optiontest:required].option[selectopt1]/type=select
# centry.cmd[optiontest:required].option[selectopt1]/envName=SELECT
# centry.cmd[optiontest:required].option[selectopt2]/type=bool
# centry.cmd[optiontest:required].option[selectopt2]/type=select
# centry.cmd[optiontest:required].option[selectopt2]/envName=SELECT
# centry.cmd[optiontest:required].option[notrequired]/required=false
optiontest:required() {
  echo "This command should not run without required options specified..."
  env | sort
}
