package main

import (
	"os"
)

// TwitterConfig env
type TwitterConfig struct {
	ConfKey     string
	ConfSecret  string
	TokenKey    string
	TokenSecret string
}

// check twice if you're uploading your precious keys and tokens to github!!!!
func conf(env TwitterConfig) TwitterConfig {

	env.ConfKey = os.Getenv("ConfKey")
	env.ConfSecret = os.Getenv("ConfKey")
	env.TokenKey = os.Getenv("TokenKey")
	env.TokenSecret = os.Getenv("TokenSecret")

	return env
}
