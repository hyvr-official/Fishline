@echo off
echo Building Go binaries

IF EXIST build (
    RMDIR /S /Q build
)

SET GORUN=true
goreleaser release --snapshot --clean

for /F "delims=" %%a in ('echo prompt $E^| cmd') do set "ESC=%%a"

echo.
echo %ESC%[32mBuild complete! Files are in the /build folder%ESC%[0m

pause