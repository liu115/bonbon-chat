# Example setup
1. Install golang on your system.
```
# apt-get install golang # Ubuntu, Debian, etc
# yum install golang     # Fedora, CentOS, etc
# pacman -S go           # ArchLinux
```

2. Set environment variable GOPATH.
```
$ export GOPATH=~/.go # "~/.go" for example. Add this line in your .*shrc
```

3. Clone git repo from hosting site.
```
$ git clone git@git.coding.net:jerry73204/bonbon.git
```

4. Initialize pre-commit GIT hook.
```
$ cd bonbon
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

# 執行伺服器

1. 安裝依賴
```
$ ./set-deps.sh # 假設在本專案跟目錄
```

2. 編譯
```
$ go install bonbon/bonbon-server
```

3. 執行
```
$ bonbon-server # 假設 $GOPATH/bin 在 $PATH 內
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
