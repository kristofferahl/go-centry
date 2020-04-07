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

# centry.cmd[test:subcommand]/description=Description for subcommand
# centry.cmd[test:subcommand]/help=Help text for sub command
test:subcommand() {
  echo "test:subcommand ($*)"
}

# centry.cmd[test:placeholder:subcommand1]/description=Description for placeholder subcommand1
# centry.cmd[test:placeholder:subcommand1]/help=Help text for placeholder subcommand1
test:placeholder:subcommand1() {
  echo "test:placeholder:subcommand1 ($*)"
}

# centry.cmd[test:placeholder:subcommand2]/description=Description for placeholder subcommand2
# centry.cmd[test:placeholder:subcommand2]/help=Help text for placeholder subcommand2
test:placeholder:subcommand2() {
  echo "test:placeholder:subcommand2 ($*)"
}
