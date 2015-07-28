## Usage
The file gives the details about the _config_ package. The purpose of this package is to manage the configurations for bonbon-server (hostname, port, database, etc).
The package does the following jobs.
* Specify default values hardcoded in source. e.g. the default port is 8080.
* Parse an ini-like config file.

The available options and default values are defined in bonbon/config/options.go. For example, you may see the default hostname and port.
```
...
var Hostname = "0.0.0.0"
...
var Port = 8080
```

Other packages can easily introduce the avaible configs. This is an example in bonbon/database/database.go.
```
import(
    ...
    "bonbon/config"
)

...
	db, err := gorm.Open(config.DatabaseDriver, config.DatabaseArgs)
```

If you would like to make an option available in the config file, simply modify bonbon/config/config.go to do this job. This requires some steps.
1. modify the struct bonbonConfig to add your entry in config file.
2. insert some sanity checking code in LoadConfigFile().
3. set the corresponding global variable(s) defined in bonbon/config/options.go in LoadConfigFile().
