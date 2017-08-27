# green 

快速创建restful接口  

* 你可以就像写本地函数一样写restful接口,接口之间甚至可以互相调用
* 采用参数注入的方式,代码更简洁
* 自动生成接口文档
* 基于[gin](https://github.com/gin-gonic/gin)二次封装

## 内置常用接口:
* 文件上传接口
* 用户注册登录接口(邮箱注册)
* 省市县地址接口
* 验证码接口
* 配置文件管理

## 使用方法

### 1. 快速开始
``` bash
# 安装 green
go get github.com/inu1255/green
# 新建项目
git clone https://github.com/inu1255/green.git your_project_dir
cd your_project_dir
# 通过main.go中的rename() 将包名都改成 your_project_dir 
go run -ldflags=-s main.go
# 修改 main.go ,取消注释 run() ,注释掉 rename()
# func main() {
#     run()
#     // 新建项目时使用该函数重命名 github.com/inu1255/green
#     // rename()
# }

# 运行
go run -ldflags=-s main.go
# 打开 http://127.0.0.1:8017/api/ 你可以看到接口文档了
```

### 2. 添加新的接口

在service文件夹下新建 hello_service.go

``` go
type HelloService struct {
    Service
}

// 所有 两个返回值且第二个返回值实现了error接口 的函数会被转化为接口
func (*HelloService)SayHello() (string,error) {
    return "hello world",nil
}

// 私有函数不会转换成接口
func (*HelloService)add(a,b int) (int,error) {
    return a+b,nil
}

// ?a=1&b=2会自动注入到参数 a,b 中
// 详细请参考 关于参数注入
func (this *HelloService)Add(a,b int) (int,error) {
    return this.add(a,b)
}

```
### 3. 自定义返回格式

* __默认返回格式__  
    函数返回值为 data,err    
    默认返回格式是json格式 {data:data,code:错误码,msg:err.Error()}   
* __自定义__  
    重写service/service.go中的Service.After函数
    

### 4. 关于参数注入

* __默认__   
    string int int64 float32 float64 会从query参数注入  (默认的使用GET方式)
    struct ptr slice map 会从post body由json格式注入 (只能有一个,默认的使用POST方式)
* __特殊的__   
    *gin.Context 会注入为当前请求的context  
    io.ReadCloser 会注入为post body  
    *multipart.FileHeader 会注入为post form file  
* __自定义__   
    可以参考 service/service.go 中的 UserManager 函数  
    并使用 maker.AddParamManager(service.UserManager) 添加注入方式  

### 5. 自动生成swagger接口文档

可以为函数写注释的方式修改接口及swagger接口文档

``` go
// 1.desc
// @desc 第一个desc接口简介
// @desc 其它的desc会显示在接口展开后
// 2.method
// 使用method修改http方法
// @method get 
// 3.path
// 使用path修改接口子路径,如果path后面什么都不接,则该函数不映射成接口
// @path /add
// 4.param
// 使用param为query参数添加介绍
// @param a 加数之一
// @param a 加数之二
// 5.tag
// 使用tag修改接口父路径
// @tag hello
func (this *HelloService)Add(a int,b int) (int,error) {
    return this.add(a,b)
}
```
函数注释信息会保存在fm.json中,以保证在没有源代码的环境中正常部署