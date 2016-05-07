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

2. 安裝 ruby 及 rake
```
$ apt-get install ruby # Ubuntu, Debian, etc
$ yum install ruby     # Fedora, CentOS, etc

$ sudo gem install rake
```

3. Clone git repo from hosting site and get into the 'bonbon' repository.
```
$ git clone git@git.coding.net:jerry73204/bonbon.git
$ cd bonbon
```

4. 若使用mysql，設定mysql
在/etc/my.cnf中加入
```
haracter-set-server=utf8
init-connect='SET NAMES utf8'
collation-server=utf8_unicode_ci
```
以確保mysql支援UTf-8

5. Check your setup (by running a build once)
```
$ rake            # Suceed if the binary 'bonbon-server' is produced
$ ./bonbon-server
```

# Usage
1. Use -h option to show usage
```
$ ./bonbon-server -h
Usage of ./bonbon-server:
...
```

2. By default, server loads the config file "bonbon.conf" in current working directory if running without any options
```
$ ./bonbon-server
```

3. Use -config option to specify the path to the config file.
```
$ ./bonbon-server -config path/to/config_file
```

# Advices
* Never directly commit to master. Always do your job on your branch.
