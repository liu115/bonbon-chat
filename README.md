# Setup (manual method)
1. Install golang on your system.
```
$ apt-get install golang # Ubuntu, Debian, etc
$ yum install golang     # Fedora, CentOS, etc
$ pacman -S go           # ArchLinux
```

2. Set environment variable GOPATH.
```
$ export GOPATH=~/.go # "~/.go" for example. Add this line in your .*shrc
```

3. Clone git repo from hosting site and get into the 'bonbon' repository.
```
$ git clone git@git.coding.net:jerry73204/bonbon.git
$ cd bonbon
```

4. Initialize pre-commit GIT hook.
```
$ ln -s ../../pre-commit.sh .git/hooks/pre-commit
```

5. Link project source in GOPATH.
```
$ ln -s $PWD ~/.go/src/bonbon # Assumed your $PWD is in the "bonbon" repo
```

6. Check your setup
```
$ go build bonbon # Suceed if no output
```

# Setup (via build automation)
1. Install golang on your system.
```
$ apt-get install golang # Ubuntu, Debian, etc
$ yum install golang     # Fedora, CentOS, etc
$ pacman -S go           # ArchLinux
```

2. Install npm and jake
```
$ apt-get install npm # Ubuntu, Debian, etc
$ yum install npm     # Fedora, CentOS, etc
$ pacman -S npm       # ArchLinux

$ npm install -g jake
```

3. Clone git repo from hosting site and get into the 'bonbon' repository.
```
$ git clone git@git.coding.net:jerry73204/bonbon.git
$ cd bonbon
```

4. Initialize pre-commit GIT hook.
```
$ ln -s ../../pre-commit.sh .git/hooks/pre-commit
```

5. Check your setup (by running a build once)
```
$ jake            # Suceed if the binary 'bonbon-server' is produced
$ ./bonbon-server
```

# 執行伺服器

1. 安裝gom
```
$ go get github.com/mattn/gom # 第一次安裝，之後可省略這步驟
```

2. 安裝依賴
```
$ gom install # 假設在本專案跟目錄
```

3. 編譯
```
$ gom build bonbon/bonbon-server
```

4. 執行
```
$ ./bonbon-server
```

# Advices
* Never directly commit to master. Always do your job on your branch.
* Master branch is updated only by "pull requests".

# Schedules
## Task # 1
* Due on 2015/07/11

### Jobs
* Register an account on coding.net and clone GIT repo (see Example Setup above)
* Finish your setup.
* jerry73204, MROS: setup GIT repo. initialize development environment
* 岳承: Clean up HTML, JS, CSS codes. Determine which frontend framework to be used.

## Task # 2
* Due on 2015/07/14

### Jobs
* Yue-Cheng
  1. remove jquery, restructure
  2. chat panel scroll bar

* jerry73204
  1. database infrastructure
  2. fixup coding style (2-space indent)

* MROS
  1. routing funcionality
  2. websocket interface

## Task # 3
* Due on 2015/07/18

### Jobs
* Yue-Cheng
  1. fix styling in chat.html

* jerry73204
  1. basic database interface (Account, ChatRoom)

* MROS
  1. websocket interface

## Break after task # 3
* Until 2015/07/21
* Hackathon Taiwan on dates 2015/07/18-19

### Summary
* Done restructure of frontend code (thanks liu115)
* building-execution (go build) and package management (gom) precedures are established
* basic database connection can be done
* basic websocket interface and a draft of websocket api are proposed (need discussion)

### Jobs
* Discussion api draft proposed by MROS (API.txt)
