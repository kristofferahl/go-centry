#!/usr/bin/env bash

up:packages () {
  echo "up:packages ($*)"
}

up:modules () {
  echo "up:modules ($*)"
}

interactive () {
  echo "interactive ($*)"
  echo -n "Enter your name: "
  read -r name
  echo "Thanks ${name:?}!"
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
