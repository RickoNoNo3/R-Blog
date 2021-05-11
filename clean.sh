#!/bin/bash

rm blog*.db >/dev/null 2>&1
rm blog*.db-* >/dev/null 2>&1
rm blog.exe >/dev/null 2>&1
rm blog.run >/dev/null 2>&1
rm ./*log.csv >/dev/null 2>&1

sqlite3op="
.read blog.sql
.quit
"
echo "$sqlite3op" | sqlite3 -batch blog.db
echo "$sqlite3op" | sqlite3 -batch blog_test.db
