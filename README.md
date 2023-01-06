# ws-im
一个轻量级的即时聊天系统, 用的是websocket. 

##本地主机临时运行命令:
```shell
$ cd deploy
$ go run ../main.go

```
进入聊天室请访问 http://127.0.0.1:8080/

## 编译命令
```shell
$ cd deploy
$ ./build_linux.sh 1.0.1

```
即可得到编译好的可执行文件`deploy/wsim`, 
注意:home.html文件必须和这个可执行文件在同一个目录.

## 或直接用docker运行
```shell
docker run -p 80:8080 cqliuz/wsim:1.0.1
```
假设docker宿主机的IP为192.168.1.113,
此时访问 http://192.168.1.113 即可进入聊天室.
