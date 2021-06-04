## windows syslog save server

### logsrv.exe
- window x64 64位操作系统
- usage 
    - ```启动 cmd```
    - ```./logsrv.exe```
### logsrv386.exe
- window 386 32位操作系统
- usage 
    - ```启动 cmd```
    - ```./logsrv386.exe```
### mock.exe 
- mock syslog  模拟syslog
- usage 
    - ```启动 cmd```
    - ```./mock.exe help```
    - ```./mock.exe -h 127.0.0.1 -p 514```


## nssm 
> nssm no-sucking service manager 安装windows服务
- 启动`cmd`
- 安装服务命令
    - ```nssm install <servicename>```
    - ```nssm install <servicename> <program>```
    - ```nssm install <servicename> <program> [<arguments>]```
- 删除服务
    - ```nssm remove```
    - ```nssm remove <servicename>```
    - ```nssm remove <servicename> confirm```
- 启动、停止服务
    - ```nssm start <servicename>```
    - ```nssm stop <servicename>```
-  查询服务状态
    - ```nssm status <servicename>```
- 服务控制命令
    - ```nssm pause <servicename>```
    - ```nssm continue <servicename>```
    - ```nssm rotate <servicename>```