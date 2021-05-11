param($killPid)

# 由于以下原因：
# 1. Windows系统没有完整的POSIX信号体系
# 2. Ctrl+C无法通过powershell进行模拟
# 3. Stop-Process是强制性的，无法触发优雅关闭
# 所以使用自定义的KillNicely类库（https://github.com/rickonono3/KillNicely）
# 由于此类库有如下特性：
# 1. 发送Ctrl+C时会给自己所在的进程里也发送出来，导致自己的进程被关闭
# 所以把调用此类库的相关逻辑写到了一个新脚本中，stop脚本开启新进程调用此脚本，以避免stop脚本被关闭
if ($null -ne $killPid) {
    $dllPath = Resolve-Path "KillNicely.dll"
    [void][reflection.assembly]::LoadFile($dllPath)
    [KillNicely.Killer]::StopProgram($killPid)
}
