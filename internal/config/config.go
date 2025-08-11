package config

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type PPlaceConfig struct {
	Port     uint           `yaml:"port"`
	Version  float32        `yaml:"version"`
	Database DatabaseConfig `yaml:"database"`
}

type Config struct {
	PPlace PPlaceConfig `yaml:"pplace"`
}
