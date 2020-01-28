#!/usr/bin/env bash

# centry.cmd.description/rotate:secrets=Rotate secrets
rotate:secrets() {
  echo "rotate:secrets ($*)"
}

# centry.cmd.description/rotate:kubernetes:workers=Rotate kubernetes worker nodes
rotate:kubernetes:workers() {
  echo "rotate:kubernetes:workers ($*)"
}

# centry.cmd.description/rotate:kubernetes:masters=Rotate kubernetes master nodes
rotate:kubernetes:masters() {
  echo "rotate:kubernetes;masters ($*)"
}
