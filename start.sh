#!/bin/bash

code=1
while [ $code -eq 1 ]; do
  ./blog.run
  code=$?
done
