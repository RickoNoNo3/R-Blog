#!/bin/bash

if [ "$1" == "silent" ]; then
  nohup env BlogRoot=. ./blog.run >/dev/null 2>&1 &
else
  env BlogRoot=. ./blog.run
fi
