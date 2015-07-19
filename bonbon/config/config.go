package config

import (
	"github.com/huandu/facebook"
)

var FBAppID = "915780538494020"
var FBAppSecret = "d7b698973175d799a27e0d129f9e19ba"
var GlobalApp = facebook.New(FBAppID, FBAppSecret)
