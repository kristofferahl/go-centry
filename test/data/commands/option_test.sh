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
# centry.cmd[optiontest:required].option[intopt]/type=integer
# centry.cmd[optiontest:required].option[intopt]/required=true
# centry.cmd[optiontest:required].option[selectopt1]/type=select
# centry.cmd[optiontest:required].option[selectopt1]/required=true
# centry.cmd[optiontest:required].option[selectopt1]/envName=SELECTOPTV1
# centry.cmd[optiontest:required].option[selectopt2]/type=select
# centry.cmd[optiontest:required].option[selectopt2]/envName=SELECTOPTV1
# centry.cmd[optiontest:required].option[selectoptv2]/type=select/v2
# centry.cmd[optiontest:required].option[selectoptv2]/envName=SELECTOPTV2
# centry.cmd[optiontest:required].option[selectoptv2]/required=true
# centry.cmd[optiontest:required].option[selectoptv2]/values=[{"name":"selectopt_v2_1"},{"name":"selectopt_v2_2"}]
# centry.cmd[optiontest:required].option[notrequired]/required=false
optiontest:required() {
  echo "This command should not run without required options specified..."
  env | sort
}
