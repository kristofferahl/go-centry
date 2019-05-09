#!/usr/bin/env bash

post() {
  echo "post ($*)"
}

put() {
  echo "put ($*)"
}

postignored() {
  echo 'should be ignored'
}

putignored() {
  echo 'should be ignored'
}
