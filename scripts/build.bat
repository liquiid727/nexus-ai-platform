@echo off
if not exist "bin" mkdir bin
go build -o bin/server.exe ./cmd/server
