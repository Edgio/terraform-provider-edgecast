:: Copyright 2021 Edgecast Inc., Licensed under the terms of the Apache 2.0 license. 
:: See LICENSE file in project root for terms.

:: This script is intended for windows environments, for Unix-like systems use the Makefile

@echo off

set HOSTNAME=github.com
set NAMESPACE=terraform-providers
set NAME=edgecast
set BINARY=terraform-provider-%NAME%.exe
set VERSION=0.5.0
set OS_ARCH=windows_amd64

go build
mkdir %APPDATA%\terraform.d\plugins\%HOSTNAME%\%NAMESPACE%\%NAME%\%VERSION%\%OS_ARCH%
move /Y %BINARY% %APPDATA%\terraform.d\plugins\%HOSTNAME%\%NAMESPACE%\%NAME%\%VERSION%\%OS_ARCH%\%BINARY%