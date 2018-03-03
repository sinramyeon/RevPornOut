package main

// TwitterConfig 설정용
type TwitterConfig struct {
	ConfKey     string
	ConfSecret  string
	TokenKey    string
	TokenSecret string
}

// 트위터 api 함수
func conf(env TwitterConfig) TwitterConfig {

	env.ConfKey = "	F5wOjWqkBPw4DXk7y4bWfALcM"
	env.ConfSecret = "vvsDzUmyKDeFIuzG0Ly13u40ezDVUSyGQo5QW3ivJDNiQd4JCE"

	env.TokenKey = "935000453225398277-rUeIivRc3CEQ3022U7uNP3DMgjfwAWE"
	env.TokenSecret = "	ypBV7vkYQd0LHqgibwAev8BtNnPgRITDkepIkkTZBfR4V"

	return env
}
