package config

type API struct {
	ApiKey           string `mapstructure:"api-key" yaml:"api-key"`
	ApiSecretKey     string `mapstructure:"api-secret-key" yaml:"api-secret-key"`
	Passphrase       string `mapstructure:"passphrase" yaml:"passphrase"`
	SimulatedTrading bool   `mapstructure:"simulated-trading" yaml:"simulated-trading,string"`
}
