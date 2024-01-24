package config

type OneApiConfig struct {
	Proxy  string `mapstructure:"proxy"`
	SToken string `mapstructure:"s-token"`
}

type PushPlus struct {
	Token string `mapstructure:"token"`
	Url   string `mapstructure:"url"`
}

type Chat struct {
	AutoPass bool `mapstructure:"autoPass"`

	Proxy          bool   `mapstructure:"proxy"`
	ProxyUrl       string `mapstructure:"proxyUrl"`
	SessionTimeOut int    `mapstructure:"sessionTimeOut"`
	Model          string `mapstructure:"model"`
}

type ServerConfig struct {
	Chat         Chat         `mapstructure:"chat"`
	OneApiConfig OneApiConfig `mapstructure:"one-api"`
	PushPlus     PushPlus     `mapstructure:"push"`
}
