Write-Debug "build"

$EXIT_ENV_ERR = 1
$EXIT_DEP_ERR = 2
$EXIT_BUILD_ERR = 3
$EXIT_RES_ERR = 4

# Check environment
Write-Output 'Checking environment'
try {
    npx -v > $null 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw
    }
} catch {
    Write-Error 'npm(npx) is not installed.'
    exit $EXIT_ENV_ERR
}
try {
    go version > $null 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw
    }
} catch {
    Write-Error 'Golang is not installed.'
    exit $EXIT_ENV_ERR
}
try {
    sqlite3 --version > $null 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw
    }
} catch {
    Write-Error 'SQLite3 is not installed.'
    exit $EXIT_ENV_ERR
}

# Install dependencies
Write-Output 'Installing dependencies'
go get -v .
if ($LASTEXITCODE -ne 0) {
    exit $EXIT_DEP_ERR
}
npm ci
if ($LASTEXITCODE -ne 0) {
    exit $EXIT_DEP_ERR
}

# Build binary
Write-Output 'Cleaning built folder'
Remove-Item built_ps1 -Force -Recurse -ErrorAction SilentlyContinue > $null
New-Item built_ps1 -ItemType Directory -ErrorAction SilentlyContinue > $null
Write-Output 'Building binary'
go build -o built_ps1/blog.exe .
if ($LASTEXITCODE -ne 0) {
    Write-Error 'Failed to build go package!'
    exit $EXIT_BUILD_ERR
}

# Init database
Write-Output 'Initing database'
$sqlite3op = "
.read blog.sql
.quit
"
Write-Output $sqlite3op | sqlite3 -batch built_ps1/blog.db
if ($LASTEXITCODE -ne 0) {
    Write-Error 'Failed to init sqlite3 database!'
    exit $EXIT_BUILD_ERR
}

# Start to minify js files
Write-Output 'Minifying *.js'
[String[]]$jsList = @()
$jsList += 'public/js/blogPage'
$jsList += 'public/js/conBg'
$jsList += 'public/js/conLogin'
$jsList += 'public/js/conNavLoc'
$jsList += 'public/js/adminConEditArticle'
$jsList += 'public/js/adminConEditDir'
$jsList += 'public/js/adminConPage'
$jsList += 'public/js/adminSettings'
$jsList += 'public/js/input-process'
foreach ($filename in $jsList) {
    npx terser -o "$filename.min.js" "$filename.js"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "> [ERROR] $filename.js"
        exit $EXIT_RES_ERR
    }
    else {
        Write-Output "> [OK] $filename.js"
    }
}

# Start to compile less files to css files
Write-Output 'Compiling *.less'
[String[]]$lessList = @()
$lessList += 'public/css/myStyles'
$lessList += 'public/css/myStyles.admin'
$lessList += 'public/css/myStyles.admin.mob'
$lessList += 'public/css/myStyles.cm'
$lessList += 'public/css/myStyles.config'
$lessList += 'public/css/myStyles.highlight'
foreach ($filename in $lessList) {
    npx lessc "$filename.less" "$filename.css"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "> [ERROR] $filename.less"
        exit $EXIT_RES_ERR
    }
    else {
        Write-Output "> [OK] $filename.less"
    }
}

# Start to minify css files (includes the less generated css files)
Write-Output 'Minifying *.css'
[String[]]$cssList = @()
$cssList += 'public/css/iconfont/iconfont'
foreach ($filename in $cssList + $lessList) {
    npx csso -i "$filename.css" -o "$filename.min.css"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "> [ERROR] $filename.css"
        exit $EXIT_RES_ERR
    }
    else {
        Write-Output "> [OK] $filename.css"
    }
}

# Copy resource files
Write-Output 'Copying resource files to built folder'
Remove-Item built_ps1/public -Force -Recurse -ErrorAction SilentlyContinue
Copy-Item public/ built_ps1/public -Force -Recurse -ErrorAction SilentlyContinue
Remove-Item built_ps1/view -Force -Recurse -ErrorAction SilentlyContinue
Copy-Item view/ built_ps1/view -Force -Recurse -ErrorAction SilentlyContinue
Copy-Item README.MD built_ps1/ -ErrorAction SilentlyContinue
Copy-Item LICENSE built_ps1/ -ErrorAction SilentlyContinue
Copy-Item config.json built_ps1/ -ErrorAction SilentlyContinue
Copy-Item start.ps1 built_ps1/ -ErrorAction SilentlyContinue
Copy-Item stop.ps1 built_ps1/ -ErrorAction SilentlyContinue
Copy-Item stop_inner.ps1 built_ps1/ -ErrorAction SilentlyContinue
Copy-Item stop.ps1 built_ps1/ -ErrorAction SilentlyContinue
Copy-Item KillNicely.dll built_ps1/ -ErrorAction SilentlyContinue
Remove-Item built_ps1/public/resource/* -Force -Recurse -ErrorAction SilentlyContinue
Remove-Item built_ps1/public/css/iconfont/.git -Force -Recurse -ErrorAction SilentlyContinue
Remove-Item built_ps1/public/js/lib/.git -Force -Recurse -ErrorAction SilentlyContinue
Remove-Item built_ps1/public/fonts/.git -Force -Recurse -ErrorAction SilentlyContinue

# Windows GUI Console

# Done
Write-Output "Done! See ./built_ps1/"
exit 0
