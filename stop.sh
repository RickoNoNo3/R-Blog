#!/bin/bash

function getPidByName() {
  pgrep blog.run | grep -E '^[0-9]+$' | awk 'NR==1{print}'
}

function getPidByPort() {
  netstat -ap | grep 'blog.run' | sed -E 's/^.*LISTEN\s*([0-9]+)\/[a-z0-9]+\s*$/\1/g' | grep -E '^[0-9]+$' | awk 'NR==1{print}'
}

pid=$(getPidByName)
if [ -n "$pid" ]; then
  # 30 seconds trying to gracefully kill
  for _ in {1..30}; do
    pid=$(getPidByName)
    if [ -z "$pid" ]; then
      break
    fi
    echo "$pid" | xargs -I {} kill {} >/dev/null 2>&1
    sleep 1
  done

  # forced kill
  pid=$(getPidByName)
  if [ -n "$pid" ]; then
    echo "$pid" | xargs -I {} kill -9 {} >/dev/null 2>&1
  fi
fi
