#!/usr/bin/env bash

# centry.cmd.description/get:env=Prints environment variables
# centry.cmd.help/get:env=Prints environment variables. Usage: ./stack get env [<...options>]
get:env() {
  env | ${SORTED}
}

# centry.cmd.description/get:files=Prints files from the current working directory
# centry.cmd.help/get:files=Prints files from the current working directory. Usage: ./stack get files [<...options>]
get:files() {
  ls -ahl | ${SORTED}
}
