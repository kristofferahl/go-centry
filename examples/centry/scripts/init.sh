#!/usr/bin/env bash

init() {
  if [[ "${DEBUG}" == true ]]; then
    set -x
  fi

  case "${SORTED}" in
    asc) SORTED='sort' ;;
    desc) SORTED='sort -r' ;;
    *) SORTED='cat' ;;
  esac

  if [[ "${NO_LOGO}" != true ]]; then
    cat <<"EOF"
                   __
  ________  ____  / /________  __
 / ___/ _ \/ __ \/ __/ ___/ / / /
/ /__/  __/ / / / /_/ /  / /_/ /
\___/\___/_/ /_/\__/_/   \__, /
                        /____/
EOF
  fi
}

init "$@"
