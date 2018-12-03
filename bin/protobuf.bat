echo off & color 0A

rem �ο����� https://github.com/google/protobuf/blob/master/cmake/README.md
rem Ĭ�ϵ�ǰ����ϵͳ�Ѱ�װ git �� cmake,�����ú��˻�������

set "WORK_DIR=%cd%"
echo %WORK_DIR%

rem ��������Ҫ��Protobuf�汾,���°汾������github�ϲ鵽 https://github.com/google/protobuf
set "PROTOBUF_VESION=v3.5.0"
echo %PROTOBUF_VESION%
set "PROTOBUF_PATH=protobuf_%PROTOBUF_VESION%"
echo %PROTOBUF_PATH%

rem ��githug����ȡprotobufԴ����
git clone -b %PROTOBUF_VESION% https://github.com/google/protobuf.git %PROTOBUF_PATH%

rem ��github����ȡgmock
cd %PROTOBUF_PATH%
git clone -b release-1.7.0 https://github.com/google/googlemock.git gmock

rem ��github����ȡgtest
cd gmock
git clone -b release-1.7.0 https://github.com/google/googletest.git gtest

cd %WORK_DIR%
rem ����VS���߼�,�൱��ָ��VS�汾,ȡ����VS�İ�װ·��
set VS_DEV_CMD="D:\Program Files (x86)\Microsoft Visual Studio 14.0\Common7\Tools\VsDevCmd.bat"
rem ���ù����ļ�������,�������ֲ�ͬ��VS�汾
set "BUILD_PATH=protobuf_%PROTOBUF_VESION%_vs2015_sln"
echo %BUILD_PATH%
if not exist %BUILD_PATH% md %BUILD_PATH%
cd %BUILD_PATH%
rem ���ñ���汾 Debug Or Release
set "MODE=Release"
echo %MODE%
if not exist %MODE% md %MODE%
cd %MODE%
echo %cd%

set "CMAKELISTS_DIR=%WORK_DIR%\%PROTOBUF_PATH%\cmake"
echo %CMAKELISTS_DIR%

rem ��ʼ�����ͱ���
call %VS_DEV_CMD%
cmake %CMAKELISTS_DIR% -G "NMake Makefiles" -DCMAKE_BUILD_TYPE=%MODE%
call extract_includes.bat
nmake /f Makefile

echo %cd%
echo %PROTOBUF_VESION%
echo %BUILD_PATH%
echo %MODE%
pause