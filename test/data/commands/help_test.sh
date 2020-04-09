#!/usr/bin/env bash

helptest() {
  echo "command args ($*)"
}

# centry.cmd[helptest:subcommand]/description=Description for subcommand
# centry.cmd[helptest:subcommand]/help=Help text for sub command
# centry.cmd[helptest:subcommand].option[opt1]/description=Help text for opt1
# centry.cmd[helptest:subcommand].option[opt1]/short=o
# centry.cmd[helptest:subcommand].option[opt1]/default=footothebar
helptest:subcommand() {
  echo "helptest:subcommand ($*)"
}

# centry.cmd[helptest:placeholder:subcommand1]/description=Description for placeholder subcommand1
# centry.cmd[helptest:placeholder:subcommand1]/help=Help text for placeholder subcommand1
# centry.cmd[helptest:placeholder:subcommand1].option[opt1]/description=Help text for opt1
helptest:placeholder:subcommand1() {
  echo "helptest:placeholder:subcommand1 ($*)"
}

# centry.cmd[helptest:placeholder:subcommand2]/description=Description for placeholder subcommand2
# centry.cmd[helptest:placeholder:subcommand2]/help=Help text for placeholder subcommand2
# centry.cmd[helptest:placeholder:subcommand2].option[opt1]/description=Help text for opt1
helptest:placeholder:subcommand2() {
  echo "helptest:placeholder:subcommand2 ($*)"
}
