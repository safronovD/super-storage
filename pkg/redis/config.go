package redis

type Config struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	HashDB   int    `yaml:"hashDB"`
	FileDB   int    `yaml:"fileDB"`
}
