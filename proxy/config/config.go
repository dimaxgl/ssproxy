package config

type Config struct {
	ListenAddress string         `yaml:"listenAddress"`
	Store         configStore    `yaml:"store"`
	Password      configPassword `yaml:"password"`
}

type configStore struct {
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}

type configPassword struct {
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}
