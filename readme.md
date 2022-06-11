                    某d wskey更换成cookie 

1 使用
```
git clone github.com/ymboom0042/wskey2Cookie
```
2 cd到根目录
```

cd path/wskey2Cookie
```
3 下载依赖
```go
go mod tidy
```
4 填写wskey到wskey.txt中
```
pin=x;wskey=x;
```
5 把自己搭建的sign配置到src/conf.yaml sign中

6 启动
```go
go run main.go
```

注：运行后会把信息会写入到根目录下的cookie.txt  
   正确：cookie: pt_pin=app_openxxx;pt_key=xxx;  
   错误：err：错误信息

巨人的肩膀：wspt.py