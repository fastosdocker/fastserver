package store

type ApplicationRsp struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
	Data Data   `json:"Data"`
}
type TemporaryList struct {
	ID          int    `json:"id"`
	CreateTime  string `json:"createTime"`
	DeletedAt   string `json:"deleted_at"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	UpdateTime  string `json:"updateTime"`
	V           string `json:"v"`
}
type List struct {
	ID          int         `json:"id"`
	Compose     interface{} `json:"compose"`
	CreateTime  string      `json:"createTime"`
	DeletedAt   string      `json:"deleted_at"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Name        string      `json:"name"`
	Note        string      `json:"note"`
	UpdateTime  string      `json:"updateTime"`
	V           string      `json:"v"`
}

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type Data struct {
	List       []List     `json:"list"`
	Pagination Pagination `json:"pagination"`
}

type ProjectConfig struct {
	// 版本
	Version string `json:"version,omitempty"`
	// 服务列表
	Services map[string]Service `json:"services,omitempty"`
	// 网络
	Networks map[string]NetworkConfig `json:"networks,omitempty"`
	// 数据卷
	Volumes map[string]VolumeConfig `json:"volumes,omitempty"`
}

type Service struct {
	// 镜像名称
	Image string `yaml:"image" json:"image,omitempty"`
	// 端口
	Ports []ServicePortConfig `yaml:"ports" json:"ports,omitempty"`
	// 重启策略
	Restart string `yaml:"restart" json:"restart,omitempty"`
	// 文件映射
	Volumes []string `yaml:"volumes" json:"volumes,omitempty"`
	// 容器名称
	ContainerName string `yaml:"container_name" json:"container_name,omitempty"`
	// 环境变量
	Environment map[string]interface{} `yaml:"environment" json:"environment,omitempty"`
	// 网络
	Networks []string `yaml:"networks" json:"networks,omitempty"`
	// 标签
	Labels map[string]string `yaml:"labels" json:"labels,omitempty"`
}

type NetworkConfig struct {
	Name       string            `yaml:"name,omitempty" json:"name,omitempty"`
	Driver     string            `yaml:"driver,omitempty" json:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty" json:"driver_opts,omitempty"`
	Ipam       IPAMConfig        `yaml:"ipam,omitempty" json:"ipam,omitempty"`
	External   External          `yaml:"external,omitempty" json:"external,omitempty"`
	Internal   bool              `yaml:"internal,omitempty" json:"internal,omitempty"`
	Attachable bool              `yaml:"attachable,omitempty" json:"attachable,omitempty"`
	Labels     Labels            `yaml:"labels,omitempty" json:"labels,omitempty"`
	EnableIPv6 bool              `yaml:"enable_ipv6,omitempty" json:"enable_ipv6,omitempty"`
}

type IPAMConfig struct {
	Driver string      `yaml:"driver,omitempty" json:"driver,omitempty"`
	Config []*IPAMPool `yaml:"config,omitempty" json:"config,omitempty"`
}

type IPAMPool struct {
	Subnet             string                 `yaml:"subnet,omitempty" json:"subnet,omitempty"`
	Gateway            string                 `yaml:"gateway,omitempty" json:"gateway,omitempty"`
	IPRange            string                 `yaml:"ip_range,omitempty" json:"ip_range,omitempty"`
	AuxiliaryAddresses map[string]string      `yaml:"aux_addresses,omitempty" json:"aux_addresses,omitempty"`
	Extensions         map[string]interface{} `yaml:",inline" json:"-"`
}

type Labels map[string]string

type External struct {
	Name     string `yaml:"name,omitempty" json:"name,omitempty"`
	External bool   `yaml:"external,omitempty" json:"external,omitempty"`
}

type VolumeConfig struct {
	Name       string            `yaml:"name,omitempty" json:"name,omitempty"`
	Driver     string            `yaml:"driver,omitempty" json:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty" json:"driver_opts,omitempty"`
	External   bool              `yaml:"external" json:"external"`
	Labels     Labels            `yaml:"labels,omitempty" json:"labels,omitempty"`
}

type ServicePortConfig struct {
	Mode      string `yaml:"mode,omitempty" json:"mode,omitempty"`
	HostIP    string `yaml:"host_ip,omitempty" json:"host_ip,omitempty"`
	Target    uint32 `yaml:"target,omitempty" json:"target,omitempty"`
	Published string `yaml:"published,omitempty" json:"published,omitempty"`
	Protocol  string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

type JSONData struct {
	Compose DockerCompose `json:"compose"`
	// CreateTime  string `json:"createTime"`
	// DeletedAt   string `json:"deleted_at"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Class       string `json:"class"`
	// UpdateTime  string `json:"updateTime"`
	V string `json:"v"`
}

type JSONDataGetId struct {
	// Compose     string `json:"compose"`
	// CreateTime  string `json:"createTime"`
	// DeletedAt   string `json:"deleted_at"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Class       string `json:"class"`
	// Note        string `json:"note"`
	// UpdateTime  string `json:"updateTime"`
	V string `json:"v"`
}
type DockerCompose struct {
	Version  string                            `yaml:"version" json:"version"`
	Services map[string]ServiceYAML            `yaml:"services" json:"services"`
	Volumes  map[string]map[string]interface{} `yaml:"volumes" json:"volumes"`
	Note     map[string]Note                   `yaml:"note" json:"note"`
}

type ServiceYAML struct {
	Image         string            `yaml:"image" json:"image"`
	ContainerName string            `yaml:"container_name" json:"container_name"`
	Restart       string            `yaml:"restart" json:"restart"`
	Volumes       []string          `yaml:"volumes" json:"volumes"`
	Ports         []string          `yaml:"ports" json:"ports"`
	Environment   map[string]string `yaml:"environment" json:"environment"`
	Privileged    string            `yaml:"privileged" json:"privileged"`
	User          string            `yaml:"user" json:"user"`
	Command       string            `yaml:"command" json:"command"`
}
type Note struct {
	Ports       []map[string]string `yaml:"ports" json:"ports"`
	Environment []map[string]string `yaml:"environment" json:"environment"`
	Volumes     []map[string]string `yaml:"volumes" json:"volumes"`
}
