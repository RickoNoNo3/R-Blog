# R-Blog

轻量级博客系统 R-Blog.

前后一体化服务器程序, HTML + Less + jQuery + Go + Sqlite 3.

-----

## 使用方法

```shell
# 自定义配置文件
cat ./config.json
# BlogRoot指定时，寻找资源相对于BlogRoot，否则相对于程序运行时的args[0]
export BlogRoot=./
# method 1
go run .
# method 2
go build -o blog.run .
./blog.run
```