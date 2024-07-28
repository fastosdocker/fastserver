package config

type conf struct {
	Port          string              // PORT 端口
	DaemonPath    string              // DaemonPath 加速器路径
	DockerHubUser []map[string]string // docker hub用户
	TLs           Tls                 `yaml:"tls"` //https证书
	Https         Https               `json:"https" yaml:"https"`
	StoreUrl      string
}

type Tls struct {
	Key string `json:"key" yaml:"key"`
	Pem string `json:"pem" yaml:"pem"`
}

type Https struct {
	Flag bool   `json:"flag" yaml:"flag"`
	Port string `json:"port" yaml:"port"`
}

type Store struct {
	Url string `json:"url" yaml:"url"`
}
