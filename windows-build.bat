@REM windres.exe -i app.rc -o res_windows_amd64.syso

@REM go build -tags=tempdll -buildmode=exe -ldflags="-s -w -H windowsgui" -o bin/config.exe .

go build -buildmode=exe -ldflags="-s -w" -o bin/task.exe .

upx bin/task.exe

@REM pause