# Fofa_scanner
*一个利用fofa api对指定ip,或者txt文本中进行批量ip的查询工具*

*主要模块*
* conver_model -- 对查询到的json数据进行解析,并且转化成csv格式
* parese_model -- 用于批量查询中，获取txt文本中的ip
* scan_model -- 调用fofa api 进行查询    //需要在apikey常量中添加自己的key
*** 
## 生成exe文件
`go build -o fofa_scan.exe`
<br>
exe文件需要cmd打开，并且保持权限足够，否则无法创建对应的文件夹，和写入文件。

1. 运行图片
   
![运行图片](./Screenshot.png)

## 更新
默认设置concurrent为10，可以并发请求






