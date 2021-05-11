#!/bin/bash

sqlite3op="
.read blog.sql
.quit
"
rm blog*.db
if ! echo "$sqlite3op" | sqlite3 -batch blog.db; then
  exit 1
fi
if ! echo "$sqlite3op" | sqlite3 -batch blog_test.db; then
  exit 1
fi
