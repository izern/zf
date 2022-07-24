- [1. 简介](#1-简介)
- [2. 安装](#2-安装)
- [3. 使用示例](#3-使用示例)
  - [3.1. append](#31-append)
  - [3.2. get](#32-get)
  - [3.3. keys](#33-keys)
  - [3.4. parse](#34-parse)
  - [3.5. set](#35-set)
  - [3.6. type](#36-type)
  - [3.7. convert](#37-convert)


## 1. 简介
zern format = zf

用于格式化各种标准格式的文件，并提供格式转换、修改指定位置的值功能。

使用jsonpath格式来处理指定位置，目前已支持

* json
* toml
* yaml

更多帮助可以查看 `zf help`

## 2. 安装

 `go install github.com/izern/zf@master`

源码安装

```bash
git clone https://github.com/izern/zf
cd zf
go install github.com/izern/zf
```

## 3. 使用示例

理论上json, toml, yaml三个子命令能提供一样的功能，以下示例仅以其中一个子命令示例。

```text

zf yaml --help

解析json格式的文本

Usage:
  zf json [flags]
  zf json [command]

Available Commands:
  append      追加值
  get         获取值
  keys        获取键列表
  parse       格式化
  set         修改值，覆盖
  type        获取指定路径值的类别

Flags:
  -h, --help   help for json

Use "zf json [command] --help" for more information about a command.
```

### 3.1. append

```bash
# 数组追加
cat test/test.yaml | zf yaml append -p .rules -v "test append"
# 字符追加
cat test/test.yaml | zf yaml append -p .port -v "1"
# 数组指定位置追加
cat test/test.yaml | zf yaml append -p .rules -i 1 -v "test append"
```

### 3.2. get

```bash
cat test/test.yaml | zf yaml get -p .rules
# 指定位置
cat test/test.yaml | zf yaml get -p .rules[1,4]
```

### 3.3. keys

```bash
cat test/test.yaml | zf yaml keys -p .proxies[0]
```

### 3.4. parse

```bash
cat test/test.yaml | zf yaml parse
```

### 3.5. set

```bash
cat test/test.yaml | zf yaml set -p .port -v 1234
```

### 3.6. type

all type 

* array
* object
* number
* string
* bool
* null

```bash
cat test/test.yaml | zf yaml type -p .port
```

### 3.7. convert

```bash
cat test/test.yaml | zf convert -f yaml -t json
```
