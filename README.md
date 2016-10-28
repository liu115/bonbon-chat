# Setup
1. 安裝go編譯器（僅確定在 go 1.6 編譯通過）
```
$ apt-get install golang # Ubuntu, Debian, etc
$ yum install golang     # Fedora, CentOS, etc
$ pacman -S go           # ArchLinux
```

2. 設定 GOPATH 環境變數及安裝gom(go 的套件管理工具)
```
$ export GOPATH=~/.go          # "~/.go" for example. Add this line in your .*shrc
$ export PATH=$PATH:$GOPATH/bin
$ go install gom
```

3. git clone repo
```
$ git clone git@git.coding.net:jerry73204/bonbon.git
$ cd bonbon
```

4. 符號連接專案到 GOPATH 下
```
$ ln -s $PWD/bonbon ~/.go/src/bonbon # Assumed your $PWD is in the "bonbon" repo
```

5. 安裝 ruby 及 rake
```
$ apt-get install ruby # Ubuntu, Debian, etc
$ yum install ruby     # Fedora, CentOS, etc
$ sudo gem install rake
```

6. 若使用mysql，設定mysql
在/etc/my.cnf中加入
```
haracter-set-server=utf8
init-connect='SET NAMES utf8'
collation-server=utf8_unicode_ci
```
以確保mysql支援UTf-8

7. 編譯前端
```
$ npm install
$ cd static
$ webpack --watch
```

8. 編譯並執行
```
$ rake                        # gom 可能會有路徑問題，請看錯誤訊息嘗試修正
$ ./bonbon-server             # 預設模式
$ ./bonbon-server -static /path/to/static -config /path/to/config
```

# Usage
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
