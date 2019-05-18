#!/usr/bin/env bash

up:packages () {
  echo "up:packages ($*)"
}

up:modules () {
  echo "up:modules ($*)"
}

up () {
  echo "up ($*)"
}

down:packages () {
  echo "down:packages ($*)"
}

down:modules () {
  echo "down:modules ($*)"
}

down () {
  echo "up ($*)"
}
