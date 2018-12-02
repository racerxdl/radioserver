# Get the version.
$commitHash = &("C:\Program Files\Git\cmd\git.exe") rev-parse HEAD
$commitHash = $commitHash.Trim()
# Write out the package.
@"
package main
//go:generate powershell .\gen.ps1
const commitHash = "$commitHash"
"@ | Out-File -Encoding ASCII -FilePath version_windows.go
# Write out the package.
@"
package main
//go:generate bash gen.sh
const commitHash = "$commitHash"
"@ | Out-File -Encoding ASCII -FilePath version_linux.go
