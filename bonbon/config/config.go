package config

import (
	"fmt"
	"strings"
	"errors"
	"code.google.com/p/gcfg"
	"github.com/huandu/facebook"
)

type bonbonConfig struct {
	Server struct {
		Hostname string
		Port     int
		Mode     string
	}

	Parameters struct {
		NumFriendsLimit   int
		ElectiveNickNames string
	}

	Database struct {
		Driver    string
		Arguments string
	}

	Facebook struct {
		AppID     string
		AppSecret string
	}
}

// LoadConfigFile load config from a file by path
func LoadConfigFile(path string) error {
	var conf bonbonConfig
	err := gcfg.ReadFileInto(&conf, path)

	if err != nil {
		return err
	}

	// sanity checks
	if conf.Server.Port < 0 || conf.Server.Port > 65535 {
		return fmt.Errorf("config: malformed server port %d", conf.Server.Port)
	}

	if conf.Server.Mode != "release" && conf.Server.Mode != "debug" && conf.Server.Mode != "test" {
		return fmt.Errorf("config: malformed server mode \"%s\", mode can be either \"release\", \"debug\" or \"test\"", conf.Server.Mode)
	}

	if conf.Parameters.NumFriendsLimit < 1 {
		return fmt.Errorf("config: malformed value NumFriendsLimit = %d. NumFriendsLimit must be positive", conf.Parameters.NumFriendsLimit)
	}

	if len(conf.Parameters.ElectiveNickNames) == 0 {
		return errors.New("config: ElectiveNickNames cannot be empry")
	}

	// populate values
	Hostname = conf.Server.Hostname
	Port     = conf.Server.Port
	Mode     = conf.Server.Mode
	Address  = fmt.Sprintf("%s:%d", Hostname, Port)

	ElectiveNickNames = strings.Split(conf.Parameters.ElectiveNickNames, ",")
	NumFriendsLimit   = conf.Parameters.NumFriendsLimit

	FBAppID     = conf.Facebook.AppID
	FBAppSecret = conf.Facebook.AppSecret
	GlobalApp   = facebook.New(FBAppID, FBAppSecret)

	DatabaseDriver = conf.Database.Driver
	DatabaseArgs   = conf.Database.Arguments

	return nil
}
