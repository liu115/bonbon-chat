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
