# 某d wskey生成cookie

1 使用
```
git clone github.com/ANMP0042/wskey2Cookie
```
2 cd到根目录
```

cd path/wskey2Cookie
```

3 填写wskey到wskey.txt中
```
pin=x;wskey=x;
```

4 把自己搭建的sign配置到src/conf.yaml sign中

5 下载依赖
```go
go mod tidy
```

6 启动
```go
go run main.go
```

注：运行后会把信息会写入到根目录下的cookie.txt  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;正确：cookie: pt_pin=app_openxxx;pt_key=xxx;  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;错误：err：错误信息

### 巨人的肩膀：wspt.py