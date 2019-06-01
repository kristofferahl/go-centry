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
}

init "$@"
