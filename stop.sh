#!/bin/bash
# TODO: 测试

function getPid() {
  netstat -ap | grep 'blog.run' | sed -E 's/^.*LISTEN\s*([0-9]+)\/[a-z0-9]+\s*$/\1/g' | grep -E '^[0-9]+$'
}

pid=$(getPid)
if [ -n "$pid" ]; then
  for _ in {1..30}; do
    echo "$pid" | xargs kill
    sleep 1
    pid=$(getPid)
    if [ -z "$pid" ]; then
      break
    fi
  done
  pid=$(getPid)
  if [ -n "$pid" ]; then
    echo "$pid" | xargs kill -9
  fi
fi
