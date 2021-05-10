param($silent)
Write-Debug "start"

if ($silent -eq "silent") {
    cmd.exe /C "start blog.exe >NUL 2>NUL"
} else {
    & .\blog.exe
}
