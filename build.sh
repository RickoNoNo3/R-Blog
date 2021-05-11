#!/bin/bash

EXIT_ENV_ERR=1
EXIT_DEP_ERR=2
EXIT_BUILD_ERR=3
EXIT_RES_ERR=4

# Check environment
echo 'Checking environment'
if ! command npx -v >&/dev/null; then
  echo 'npm(npx) is not installed.' >&2
  exit $EXIT_ENV_ERR
fi
if ! command go version >&/dev/null; then
  echo 'Golang is not installed.' >&2
  exit $EXIT_ENV_ERR
fi
if ! command sqlite3 --version >&/dev/null; then
  echo 'SQLite3 is not installed.' >&2
  exit $EXIT_ENV_ERR
fi

# Install dependencies
echo 'Installing dependencies'
if ! go get -v .; then
  exit $EXIT_DEP_ERR
fi
if ! npm clean-install; then
  exit $EXIT_DEP_ERR
fi

# Build binary
echo 'Cleaning built folder'
rm -rf built_sh
mkdir built_sh
echo 'Building binary'
if ! go build -o built_sh/blog.run .; then
  echo 'Failed to build go package!' >&2
  exit $EXIT_BUILD_ERR
fi

# Init database
echo 'Initing database'
sqlite3op="
.read blog.sql
.quit
"
if ! echo "$sqlite3op" | sqlite3 -batch built_sh/blog.db; then
  echo 'Failed to init sqlite3 database!' >&2
  exit $EXIT_BUILD_ERR
fi

# Start to minify js files
echo 'Minifying *.js'
jsList[0]='public/js/blogPage'
jsList[1]='public/js/conBg'
jsList[2]='public/js/conLogin'
jsList[3]='public/js/conNavLoc'
jsList[4]='public/js/adminConEditArticle'
jsList[5]='public/js/adminConEditDir'
jsList[6]='public/js/adminConPage'
jsList[7]='public/js/adminSettings'
jsList[8]='public/js/input-process'
for filename in ${jsList[*]}; do
  if ! npx terser <"$filename.js" >"$filename.min.js"; then
    echo "> [ERROR] $filename.js"
    exit $EXIT_RES_ERR
  else
    echo "> [OK] $filename.js"
  fi
done

# Start to compile less files to css files
echo 'Compiling *.less'
lessList[0]='public/css/myStyles'
lessList[1]='public/css/myStyles.admin'
lessList[2]='public/css/myStyles.admin.mob'
lessList[3]='public/css/myStyles.cm'
lessList[4]='public/css/myStyles.config'
lessList[5]='public/css/myStyles.highlight'
for filename in ${lessList[*]}; do
  if ! npx lessc "$filename.less" "$filename.css"; then
    echo "> [ERROR] $filename.less"
    exit $EXIT_RES_ERR
  else
    echo "> [OK] $filename.less"
  fi
done

# Start to minify css files (includes the less generated css files)
echo 'Minifying *.css'
cssList[0]='public/css/iconfont/iconfont'
cssListAll=("${cssList[*]}" "${lessList[*]}")
for filename in ${cssListAll[*]}; do
  if ! npx csso <"$filename.css" >"$filename.min.css"; then
    echo "> [ERROR] $filename.css"
    exit $EXIT_RES_ERR
  else
    echo "> [OK] $filename.css"
  fi
done

# Copy resource files
echo 'Copying resource files to built folder'
rm -rf built_sh/public &>/dev/null
cp -r public/ built_sh/public
rm -rf built_sh/view &>/dev/null
cp -r view/ built_sh/view
cp README.MD built_sh/
cp LICENSE built_sh/
cp config.json built_sh/
cp start.sh built_sh/
cp stop.sh built_sh/
rm -rf built_sh/public/resource/* &>/dev/null
rm -rf built_sh/public/css/iconfont/.git &>/dev/null
rm -rf built_sh/public/js/lib/.git &>/dev/null
rm -rf built_sh/public/fonts/.git &>/dev/null

# Done
echo "Done! See ./built_sh/"
exit 0
