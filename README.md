# go-server
goland 游戏服务器架构，mmo架构，ai行为树，配置data

设计之初，建立在actor模式下的；rpc，以及消息驱动，rpc无需注册，支持int，数据，struct（struct必须要注册结构题即可）；sql封装简单的orm，具体看demo

websocket模式下，要在net，websocket注视掉如下代码：https://studygolang.com/articles/14842

1.配置golang的gopath和goroot

2.bin目录下的sxz_server.cfg配置数据库以及端口

3.数据库在sql文件目录下生产

4.win下执行build.bat,start.bat

5.linux下执行build.sh,start.sh

