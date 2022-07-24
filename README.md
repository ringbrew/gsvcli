# gsvcli

## How to install

go install -v github.com/ringbrew/gsvcli/gsv@latest

## How to use

### gsv init github.com/rinbgbrew/demo

初始化脚手架项目

### gsv gen domain {domainName}

生成domain文件夹内容

### gsv gen grpc -I xxx -P proto

参数说明：

-P 项目内的proto路径

-I proto依赖执行路径

根据proto文件内的内容生成grpc service，proto这个参数可以不传，默认为proto

### gsv install

安装gsv依赖
