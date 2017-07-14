# Setup
1. 安裝go編譯器（僅確定在 go 1.6 編譯通過）
``` sh
$ apt-get install golang # Ubuntu, Debian, etc
$ yum install golang     # Fedora, CentOS, etc
$ pacman -S go           # ArchLinux
```

2. 設定 GOPATH 環境變數及安裝gom(go 的套件管理工具)
``` sh
$ export GOPATH=~/.go          # "~/.go" for example. Add this line in your .*shrc
$ export PATH=$PATH:$GOPATH/bin
```

3. 下載 git repo
``` sh
$ git clone git@git.coding.net:jerry73204/bonbon.git
$ cd bonbon
```

4. 符號連接專案到 GOPATH 下
``` sh
$ ln -s $PWD/bonbon ~/.go/src/bonbon # Assumed your $PWD is in the "bonbon" repo
```

5. 若使用mysql，設定mysql
在/etc/my.cnf中加入
```
haracter-set-server=utf8
init-connect='SET NAMES utf8'
collation-server=utf8_unicode_ci
```
以確保mysql支援UTf-8

6. 編譯前端
``` sh
$ npm install
$ cd static
$ webpack --watch
# 可給予環境變數來控制行爲
# silent=true|false （預設 false ）控制是否在前端打印訊息
# NODE_ENV=development|production （預設 development） production 會關閉打印、並醜化JavaScript
$ env NODE_ENV=production webpack # 佈屬時使用
```

7. 編譯並執行

注意！目前直接採用 go get 控制第三方函式庫，所以，如果所在機器 go 函式庫版本複雜，可能會出問題

``` sh
$ sh install.sh
$ go build bonbon/bonbon-server
$ ./bonbon-server             # 預設模式
$ ./bonbon-server -static /path/to/static -config /path/to/config
```

# 使用
1. 用 -h 來查看幫助
```
$ ./bonbon-server -h
```

2. 預設會抓取同目錄下的 bonbon-develop.conf 作為設定檔，static/ 作為靜態檔案位置
```
$ ./bonbon-server
```

3. 使用 -config 參數指定設定檔， -static 來指定靜態檔案的位置
```
$ ./bonbon-server -config path/to/config_file -static path/to/static/
```
