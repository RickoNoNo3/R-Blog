Write-Debug "stop"


$blog = Get-Process "blog" -ea SilentlyContinue
for ($i = 0; $i -lt 30; $i++) {
    if ($null -ne $blog) {
        powershell .\stop_inner.ps1 $blog.Id
        Start-Sleep 1
    } else {
        break
    }
    $blog = Get-Process "blog" -ea SilentlyContinue
}
if ($null -ne $blog) {
    Stop-Process $blog -Force
    Remove-Variable $blog
}
