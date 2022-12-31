package config

type Config struct {
}

func New(envDir string) (Config, error) {
	cfg := Config{}
	err := Process(envDir, &cfg)
	return cfg, err
}

func Process(envDir string, cfg *Config) error {
	return nil
}
