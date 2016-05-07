# Setup
1. 安裝go編譯器（目前只在 go 1.5 編譯成功）
```
$ apt-get install golang # Ubuntu, Debian, etc
$ yum install golang     # Fedora, CentOS, etc
$ pacman -S go           # ArchLinux
```

2. 設定 GOPATH 環境變數
```
$ export GOPATH=~/.go          # "~/.go" for example. Add this line in your .*shrc
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

7. 編譯並執行
```
$ rake
$ ./bonbon-server             # 預設模式
$ ./bonbon-server -static /path/to/static -config /path/to/config
```

# Usage
1. 用 -h 來查看幫助
```
$ ./bonbon-server -h
...
```

2. 預設會抓取同目錄下的 bonbon-develop.conf 作為設定檔，static/ 作為靜態檔案位置
```
$ ./bonbon-server
```

3. 使用 -config 參數指定設定檔， -static 來指定靜態檔案的位置
```
$ ./bonbon-server -config path/to/config_file -static path/to/static/
```
