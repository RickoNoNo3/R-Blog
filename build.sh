#!/bin/bash

# Cd to the owner dir of this shell
loc=$(echo "$0" | sed 's/^\(.*\/\)[^/]*$/\1/')
if [ "$loc" != "./" ]; then
  cd "$loc" || exit 127
fi

# Check environment
echo 'Checking environment'
if ! command npx -v >&/dev/null; then
  echo 'npm(npx) is not installed.' >&2
  exit 1
fi
if ! command go version >&/dev/null; then
  echo 'Golang is not installed.' >&2
  exit 1
fi

# Install dependencies
echo 'Installing dependencies'
go get .
npm install

# Build binary
echo 'Cleaning built folder'
rm -rf built
mkdir built
echo 'Building binary'
go build -o built/blog.run .
if [ "$0" != "0" ]; then
  echo 'Failed to build go package!' >&2
  exit 2
fi

# Start to minify js files
echo 'Minifying *.js'
MinifyJsList[0]='public/js/blogPage'
MinifyJsList[1]='public/js/conBg'
MinifyJsList[2]='public/js/conLogin'
MinifyJsList[3]='public/js/conNavLoc'
for filename in ${MinifyJsList[*]}; do
  npx terser <"$filename.js" >"$filename.min.js"
  if [ "$0" == "0" ]; then
    echo "> [ERROR] $filename"
  else
    echo "> [OK] $filename"
  fi
done

# Start to minify css files
echo 'Minifying *.css'
MinifyCssList[0]='public/css/myStyles'
MinifyCssList[1]='public/css/myStyles.mob'
MinifyCssList[2]='public/css/myStyles.highlight'
for filename in ${MinifyCssList[*]}; do
  npx csso <"$filename.css" >"$filename.min.css"
  if [ "$0" == "0" ]; then
    echo "> [ERROR] $filename"
  else
    echo "> [OK] $filename"
  fi
done

# TODO: Less

# Copy resource files
echo 'Copying resource files to built folder'
rm -rf built/public &> /dev/null
cp -r public/ built/public
rm -rf built/view &> /dev/null
cp -r view/ built/view
cp README.MD built/
cp LICENSE built/
cp blog.db built/

# Done
echo "Done! See $loc/built/"
