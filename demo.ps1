# Windows PowerShell script to run the demo
Write-Host "Starting the hot reload demonstration..." -ForegroundColor Green
go run . --root ./testserver --build "go build -o ./bin/server.exe ./testserver" --exec "./bin/server.exe"