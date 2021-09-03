#!/usr/bin/env bash

commandtest() {
  echo "command args ($*)"
}

commandtest:subcommand() {
  echo "subcommand args ($*)"
}

# centry.cmd[commandtest:options:args].option[cmdstringopt]/type=string
# centry.cmd[commandtest:options:args].option[cmdboolopt]/type=bool
# centry.cmd[commandtest:options:args].option[cmdsel1]/type=select
# centry.cmd[commandtest:options:args].option[cmdsel1]/envName=CMDSELECTOPT
# centry.cmd[commandtest:options:args].option[cmdsel2]/type=select
# centry.cmd[commandtest:options:args].option[cmdsel2]/envName=CMDSELECTOPT
commandtest:options:args() {
  echo "command args ($*)"
}

# centry.cmd[commandtest:options:printenv].option[cmdstringopt]/type=string
# centry.cmd[commandtest:options:printenv].option[cmdboolopt]/type=bool
# centry.cmd[commandtest:options:printenv].option[cmdsel1]/type=select
# centry.cmd[commandtest:options:printenv].option[cmdsel1]/envName=CMDSELECTOPT
# centry.cmd[commandtest:options:printenv].option[cmdsel2]/type=select
# centry.cmd[commandtest:options:printenv].option[cmdsel2]/envName=CMDSELECTOPT
# centry.cmd[commandtest:options:printenv].option[dashed-opt]/type=string
commandtest:options:printenv() {
  env | sort
}

commandtest:exitcode() {
  exit 111
}
