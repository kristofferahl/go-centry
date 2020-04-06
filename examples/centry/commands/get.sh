#!/usr/bin/env bash

# centry.cmd[get:env]/description=Prints environment variables
# centry.cmd[get:env]/help=Prints environment variables. Usage: ./stack get env [<...options>]
# centry.cmd[get:env].option[filter]/short=f
# centry.cmd[get:env].option[filter]/description=Filters environment variables based on the provided value
# centry.cmd[get:env].option[sanitize]/type=bool
# centry.cmd[get:env].option[sanitize]/description=Clean output so that no secrets are leaked
# centry.cmd[get:env].option[sanitize]/envName=SANITIZE_OUTPUT
get:env() {
  local output
  output="$(env | ${SORTED} | grep "${FILTER}")"

  if [[ ${SANITIZE_OUTPUT} == true ]]; then
    echo "${output}" | sed 's/\=.*$/=***/'
  else
    echo "${output}"
  fi
}

# centry.cmd[get:files]/description=Prints files from the current working directory
# centry.cmd[get:files]/help=Prints files from the current working directory. Usage: ./stack get files [<...options>]
get:files() {
  ls -ahl | ${SORTED}
}
