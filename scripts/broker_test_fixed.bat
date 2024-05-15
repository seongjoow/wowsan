@echo off
setlocal enabledelayedexpansion

SET "TICK_DIR_NAME=./log/tickLogger/"
SET "HOP_DIR_NAME=./log/hopLogger/"
SET "MAX_INDEX=0"

:: 디렉토리 존재 여부 확인
IF NOT EXIST "%TICK_DIR_NAME%" (
    echo Directory %TICK_DIR_NAME% not found.
    exit /b 1
)
IF NOT EXIST "%HOP_DIR_NAME%" (
    echo Directory %HOP_DIR_NAME% not found.
    exit /b 1
)

:: 가장 큰 인덱스 찾기
:loop
IF EXIST "./%TICK_DIR_NAME%%MAX_INDEX%" (
    SET /A MAX_INDEX=%MAX_INDEX% + 1
    echo %MAX_INDEX%
    GOTO loop
)

FOR /L %%P IN (50001, 1, 50010) DO (
    START "Broker Service on Port %%P" cmd /k go run .\cmd\broker\main_test_fixed.go --port %%P --dir_index %MAX_INDEX%
)

endlocal