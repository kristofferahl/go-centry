#!/usr/bin/env bash

# centry.cmd[rotate:secrets]/description=Rotate secrets
rotate:secrets() {
  echo "rotate:secrets ($*)"
}

# centry.cmd[rotate:kubernetes:workers]/description=Rotate kubernetes worker nodes
rotate:kubernetes:workers() {
  echo "rotate:kubernetes:workers ($*)"
}

# centry.cmd[rotate:kubernetes:masters]/description=Rotate kubernetes master nodes
rotate:kubernetes:masters() {
  echo "rotate:kubernetes;masters ($*)"
}
