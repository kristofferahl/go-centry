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
  output="$(env | ${SORTED:-cat})"
  [[ -n "${FILTER}" ]] && output="$(echo "${output:-}" | grep "${FILTER}")"

  if [[ ${SANITIZE_OUTPUT} == true ]]; then
    echo "${output}" | sed 's/\=.*$/=***/'
  else
    echo "${output}"
  fi
}

# centry.cmd[get:files]/description=Prints files from the current working directory
# centry.cmd[get:files]/help=Prints files from the current working directory. Usage: ./stack get files [<...options>]
# centry.cmd[get:files].option[hidden]/description=A hidden option
# centry.cmd[get:files].option[hidden]/hidden=true
get:files() {
  ls -ahl | ${SORTED}
}

# centry.cmd[get:hidden]/hidden=true
get:hidden() {
  echo "This subcommand won't be displayed in help output"
}

# centry.cmd[get:error]/description=A command that will always generate an error
get:error() {
  echo "I will generate an error"
  exit 123
}

# centry.cmd[get:required].option[abc]/required=true
# centry.cmd[get:required].option[def]/required=false
get:required() {
  echo "This subcommand has a required option"
  env | sort
}

# centry.cmd[get:selected].option[abc]/type=select
# centry.cmd[get:selected].option[abc]/envName=SELECTED
# centry.cmd[get:selected].option[abc]/required=true
# centry.cmd[get:selected].option[def]/type=select
# centry.cmd[get:selected].option[def]/envName=SELECTED
# centry.cmd[get:selected].option[def]/required=true
get:selected() {
  echo "The selected value was ${SELECTED:?} (select v1)"
}

# centry.cmd[get:selectedv2].option[selected]/type=select/v2
# centry.cmd[get:selectedv2].option[selected]/envName=SELECTED
# centry.cmd[get:selectedv2].option[selected]/required=true
# centry.cmd[get:selectedv2].option[selected]/values=[{"name":"abc","short":"a","value":"val1"},{"name":"def","short":"d","value":"val2"}]
get:selectedv2() {
  echo "The selected value was ${SELECTED:?} (select v2)"
}

# centry.cmd[get:url]/description=Call a url
# centry.cmd[get:url].option[url]/description=URL to call
# centry.cmd[get:url].option[url]/required=true
# centry.cmd[get:url].option[max-retries]/description=Maximum number of retries
# centry.cmd[get:url].option[max-retries]/type=integer
# centry.cmd[get:url].option[max-retries]/default=3
get:url() {
  echo "Calling ${URL:?} a maximum of ${MAX_RETRIES:?} time(s)"
  echo

  local success=false
  local attempts=0
  until [[ ${success:?} == true ]]; do
    ((attempts++))

    if ! curl "${URL:?}"; then
      if [[ ${attempts:?} -ge ${MAX_RETRIES:?} ]]; then
        echo
        echo "Max retries reached... exiting!"
        return 1
      fi
      sleep 1
    else
      success=true
    fi
  done
}
