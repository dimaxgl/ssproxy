package config

type Config struct {
	ListenAddress string      `yaml:"listenAddress"`
	Store         configStore `yaml:"store"`
}

type configStore struct {
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}
