# Lambda

## 安装 aws cli

```bash
sudo snap install aws-cli --classic
```

## 配置 aws config

访问[我的安全凭证](https://console.amazonaws.cn/iam/home?region=cn-north-1#/security_credentials)，创建一个新的访问密钥。

``` bash
aws configure
AWS Access Key ID [None]: YOUR_ACCESS_KEY
AWS Secret Access Key [None]: YOUR_SECRET_KEY
Default region name [None]: YOUR_REGION
Default output format [None]: None
```

## 为 Lambda 创建 Role

### 创建 Role 配置文件

新建一个 `json` 格式的文件，比如叫 `create-lambda-role.json`，内容为

``` json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "lambda.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```

### 创建 Role

通过 `aws` 命令创建 `Role` 

``` bash
aws iam create-role \
    --role-name lambda-books-executor \
    --assume-role-policy-document file:///$(pwd)/create-lambda-role.json
```

创建成功的输出如下

``` json
{
    "Role": {
        "RoleId": "AROA5ZNYJE74HJ3O4YYRZ",
        "RoleName": "lambda-books-executor",
        "Arn": "arn:aws-cn:iam::${YOUR_ACCOUNT}:role/lambda-books-executor",
        "Path": "/",
        "AssumeRolePolicyDocument": {
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                    },
                    "Action": "sts:AssumeRole"
                }
            ],
            "Version": "2012-10-17"
        },
        "CreateDate": "2020-06-22T02:19:58Z"
    }
}
```

## 编写 Lambda

示例代码

``` go
package main

import (
    "fmt"
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
    Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
    return fmt.Sprintf("Hello %s!", name.Name ), nil
}

func main() {
    lambda.Start(HandleRequest)
}
```

+ **package main**：在 Go 中，包含 func main() 的程序包必须始终名为 main。
+ **import**：请使用此包含您的 Lambda 函数需要的库。在此实例中，它包括：
    - **上下文**：[Go 中的 AWS Lambda 上下文对象](https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/golang-context.html)。
    - **fmt**：用于格式化您的函数返回的值的 Go [格式化](https://golang.org/pkg/fmt/)对象。
    - **github.com/aws/aws-lambda-go/lambda**：如前所述，实现适用于 `Go` 的 `Lambda` 编程模型。
+ **func HandleRequest(ctx context.Context, name MyEvent) (string, error)**：这是您的 Lambda 处理程序签名且包括将执行的代码。此外，包含的参数表示以下含义：
    - **ctx context.Context**：为您的 `Lambda` 函数调用提供运行时信息。`ctx` 是您声明的变量，用于利用通过 [Go 中的 AWS Lambda 上下文对象](https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/golang-context.html) 提供的信息。
    - **name MyEvent**：变量名称为 name 的输入类型，其值将在 `return` 语句中返回。
    - **string, error**：返回两个值：成功时的字符串和标准[错误](https://golang.org/pkg/builtin/#error)信息。有关自定义错误处理的更多信息，请参阅[Go 中的 AWS Lambda 函数错误](https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/golang-exceptions.html)。
    - **return fmt.Sprintf("Hello %s!", name), nil**：只返回格式化“Hello”问候语和您在输入事件中提供的姓名。`nil` 表示没有错误，函数已成功执行。
+ **func main()**：执行您的 `Lambda` 函数代码的入口点。该项为必填项。
    - 通过在 `func main(){}` 代码括号之间添加 `lambda.Start(HandleRequest)`，您的 `Lambda` 函数将会执行。

### 使用结构化类型的 Lambda 函数处理程序

在上述示例中，输入类型是简单的字符串。但是，您也可以将结构化事件传递到您的函数处理程序：

``` go
package main
 
import (
        "fmt"
        "github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
        Name string `json:"What is your name?"`
        Age int     `json:"How old are you?"`
}
 
type MyResponse struct {
        Message string `json:"Answer:"`
}
 
func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
        return MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
}
 
func main() {
        lambda.Start(HandleLambdaEvent)
}
```

然后，您的请求将如下所示：

``` json
# request
{
    "What is your name?": "Jim",
    "How old are you?": 33
}
```

而响应将如下所示：

``` json
# response
{
    "Answer": "Jim is 33 years old!"
}
```

若要导出，事件结构中的字段名称必须大写。有关来自 AWS 事件源的处理事件的更多信息，请参见 [aws-lambda-go/events](https://github.com/aws/aws-lambda-go/tree/master/events)。

#### 有效处理程序签名

在 `Go` 中构建 `Lambda` 函数处理程序时，您有多个选项，但您必须遵守以下规则：
* 处理程序必须为函数。
* 处理程序可能需要 0 到 2 个参数。如果有两个参数，则第一个参数必须实现 `context.Context`。
* 处理程序可能返回 0 到 2 个参数。如果有一个返回值，则它必须实现 `error`。如果有两个返回值，则第二个值必须实现 `error`。有关实现错误处理信息的更多信息，请参阅[Go 中的 AWS Lambda 函数错误](https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/golang-exceptions.html)。

下面列出了有效的处理程序签名。`TIn` 和 `TOut` 表示类型与 `encoding/json` 标准库兼容。有关更多信息，请参阅 [func Unmarshal](https://golang.org/pkg/encoding/json/#Unmarshal)，以了解如何反序列化这些类型。

``` go
func ()
func () error
func (TIn), error
func () (TOut, error)
func (context.Context) error
func (context.Context, TIn) error
func (context.Context) (TOut, error)
func (context.Context, TIn) (TOut, error)
```

### 使用全局状态

您可以声明并修改独立于 `Lambda` 函数的处理程序代码的全局变量。此外，您的处理程序可能声明一个 `init` 函数，该函数在加载您的处理程序时执行。这在 `AWS Lambda` 中行为方式相同，正如在标准 `Go` 程序中一样。您的 `Lambda` 函数的单个实例将不会同时处理多个事件。

``` go
package main
 
import (
        "log"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/s3"
        "github.com/aws/aws-sdk-go/aws"
)
 
var invokeCount = 0
var myObjects []*s3.Object
func init() {
        svc := s3.New(session.New())
        input := &s3.ListObjectsV2Input{
                Bucket: aws.String("examplebucket"),
        }
        result, _ := svc.ListObjectsV2(input)
        myObjects = result.Contents
}
 
func LambdaHandler() (int, error) {
        invokeCount = invokeCount + 1
        log.Print(myObjects)
        return invokeCount, nil
}
 
func main() {
        lambda.Start(LambdaHandler)
}
```

## 部署 Lambda

``` bash
aws lambda create-function \
    --function-name first-try \
    --runtime go1.x \
    --zip-file fileb://function.zip \
    --handler main \
    --role arn:aws-cn:iam::${YOUR_ACCOUNT}:role/lambda-books-executor
```

## 参考文档

[使用 Go 构建 Lambda 函数](https://docs.aws.amazon.com/zh_cn/lambda/latest/dg/lambda-golang.html)

[[译] 使用 Go 和 AWS Lambda 构建无服务 API](https://juejin.im/post/5af4082f518825672a02f262)