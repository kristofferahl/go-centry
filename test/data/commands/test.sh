#!/usr/bin/env bash

# centry.cmd[test:args].option[cmdstringopt]/type=string
# centry.cmd[test:args].option[cmdboolopt]/type=bool
# centry.cmd[test:args].option[cmdsel1]/type=select
# centry.cmd[test:args].option[cmdsel1]/envName=CMDSELECTOPT
# centry.cmd[test:args].option[cmdsel2]/type=select
# centry.cmd[test:args].option[cmdsel2]/envName=CMDSELECTOPT
test:args() {
  echo "test:args ($*)"
}

# centry.cmd[test:env].option[cmdstringopt]/type=string
# centry.cmd[test:env].option[cmdboolopt]/type=bool
# centry.cmd[test:env].option[cmdsel1]/type=select
# centry.cmd[test:env].option[cmdsel1]/envName=CMDSELECTOPT
# centry.cmd[test:env].option[cmdsel2]/type=select
# centry.cmd[test:env].option[cmdsel2]/envName=CMDSELECTOPT
test:env() {
  env | sort
}
