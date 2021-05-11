Write-Output "clean"

Remove-Item blog*.db -Force -Recurse > $null 2>&1
Remove-Item blog*.db-* -Force -Recurse > $null 2>&1
Remove-Item blog.exe -Force -Recurse > $null 2>&1
Remove-Item blog.run -Force -Recurse > $null 2>&1
Remove-Item ./*log.csv -Force -Recurse > $null 2>&1

$sqlite3op="
.read blog.sql
.quit
"
Write-Output $sqlite3op | sqlite3 -batch blog.db
Write-Output $sqlite3op | sqlite3 -batch blog_test.db
