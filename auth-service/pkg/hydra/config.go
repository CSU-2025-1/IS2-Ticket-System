package hydra

type Config struct {
	AdminURL string `yaml:"admin_url"`
	Debug    bool   `yaml:"debug"`
}
