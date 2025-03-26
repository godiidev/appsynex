@echo off
setlocal

if "%1" == "build" (
    echo Building application...
    if not exist bin mkdir bin
    go build -o bin\appsynex-api.exe .\cmd\api
) else if "%1" == "run" (
    go run .\cmd\api\main.go
) else if "%1" == "test" (
    go test -v .\...
) else if "%1" == "clean" (
    echo Cleaning build files...
    if exist bin rmdir /s /q bin
    go clean
) else if "%1" == "seed" (
    echo Running seed data...
    go run .\scripts\seed.go
) else (
    echo Unknown command: %1
    echo Available commands: build, run, test, clean, seed
)

endlocal