chcp 65001
Write-Debug "stop"

# getPid
$blog = Get-Process "blog" -ea SilentlyContinue

if ($null -ne $blog) {
    # 30 seconds trying to gracefully kill
    for ($i = 0; $i -lt 30; $i++) {
        $blog = Get-Process "blog" -ea SilentlyContinue
        if ($null -eq $blog) {
            break
        }
        powershell -File .\stop_inner.ps1 $blog.Id
        Start-Sleep -Seconds 1
    }
    # forced kill
    $blog = Get-Process "blog" -ea SilentlyContinue
    if ($null -ne $blog) {
        Stop-Process $blog -Force -ea SilentlyContinue
        Remove-Variable $blog -ea SilentlyContinue
    }
}
