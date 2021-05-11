chcp 65001
Write-Debug "start"

$code = 1
while ($code -eq 1) {
    .\blog.exe
    $code = $LASTEXITCODE
}
