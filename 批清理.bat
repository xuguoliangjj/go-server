@echo off
echo �������obj pch idb pdb ncb opt plg res sbr ilk suo�ļ������Ե�......
pause
del /f /s /q .\*.obj
del /f /s /q .\*.pch
del /f /s /q .\*.idb
del /f /s /q .\*.pdb
del /f /s /q .\*.ncb 
del /f /s /q .\*.opt 
del /f /s /q .\*.plg
del /f /s /q .\*.sdf
del /f /s /q .\*.sbr
del /f /s /q .\*.ilk
del /f /s /q .\*.aps
del /f /s /q .\*.ipch
del /f /s /q .\*.dmp
del /f /s /q .\*.log
del /f /s /q .\*.err
del /f /s /q .\*.DS_Store
del /f /s /q server.exe
del /f /s /q client.exe
del /f /s /q nohup.out
rd  /s /q .\pkg
rd  /s /q .\.idea



echo ����ļ���ɣ�
echo. & pause