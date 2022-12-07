package config

type Config struct {
	Width  int
	Height int
}

var C *Config

func init() {
	C = &Config{
		Width:  640,
		Height: 360,
	}
}
