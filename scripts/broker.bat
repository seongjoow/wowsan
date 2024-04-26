@echo off
SETLOCAL EnableDelayedExpansion

FOR /L %%P IN (50001, 1, 50010) DO (
    START "Broker Service on Port %%P" cmd /k go run .\cmd\broker\main.go --port %%P
)
