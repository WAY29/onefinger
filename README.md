# onefinger
一个简单的指纹识别工具，规则提取自goby并稍作修改，某个知识星球的作业
## Usage
```
Usage: onefinger [-v] | (-t=<target>... | --tf=<targetFile>) [--threads=<threads>] [--timeout=<timeout>]

Simple website fingerprinting tool

Options:
  -t, --target    Target url
      --tf        Target url file
      --threads   Thread number (default 10)
      --timeout   Request timeout (default 20)
  -v, --version   Show the version and exit
```
## build
```
go build main.go
```
or in windows
```cmd
cd releases; .\build.bat
```
## example
```
onefinger -t http://106.12.46.49/ --threads 10
onefinger --tf ./targets.txt --threads 10
```