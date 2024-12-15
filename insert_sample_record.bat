@echo off

:: 시스템 환경변수 설정
set "PSQL_PATH=C:\Program Files\PostgreSQL\17\bin"

:: 환경변수에 경로가 이미 설정되어 있는지 확인
setlocal enableextensions enabledelayedexpansion
set "PATH_CHECK=!PATH:%PSQL_PATH%=!"
if "!PATH_CHECK!" == "%PATH%" (
    echo PostgreSQL path already exists.
) else (
    echo Add PostgreSQL path.
    setx PATH "%PATH%;%PSQL_PATH%"
)

psql -U postgres -d test -h localhost -p 5432 -f insert_data.sql
