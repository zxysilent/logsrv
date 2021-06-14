rsrc -manifest main.manifest -ico favicon.ico -o rsrc.syso
go build -ldflags="-s -w -H=windowsgui"
@REM upx -9 ./vmgui.exe
