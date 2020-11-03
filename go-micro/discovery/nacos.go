package discovery

import (
	"strconv"

	"github.com/lei006/nacos-sdk-go/clients"
	"github.com/lei006/nacos-sdk-go/clients/config_client"
	"github.com/lei006/nacos-sdk-go/clients/naming_client"
	"github.com/lei006/nacos-sdk-go/common/constant"
	"github.com/lei006/nacos-sdk-go/model"
	"github.com/lei006/nacos-sdk-go/vo"
	"github.com/tidwall/gjson"
)

type DiscoveryContent struct {
	ServerConfigList []constant.ServerConfig
	ClientConfig     constant.ClientConfig

	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient
}

type Service struct {
	Ip          string
	Port        uint64
	Name        string
	Weight      float64
	Enable      bool
	Healthy     bool
	Ephemeral   bool
	Metadata    map[string]string
	ClusterName string // default value is DEFAULT
	GroupName   string // default value is DEFAULT_GROUP
	Ins         *model.Instance
}

type ListenerOption func(group, dataId, data string)

var instance DiscoveryContent

func GetInstance() *DiscoveryContent {
	return &instance
}

func (this *DiscoveryContent) Run(server, client string) error {

	// 取得服务器配置
	{
		configNacosServer := server
		serverConfigs := gjson.Get(configNacosServer, "server").Array()

		for _, configItem := range serverConfigs {
			var itemCfg = constant.ServerConfig{}
			strItem := configItem.String()

			itemCfg.IpAddr = gjson.Get(strItem, "IpAddr").String()
			itemCfg.Port = gjson.Get(strItem, "Port").Uint()
			itemCfg.ContextPath = gjson.Get(strItem, "ContextPath").String()
			this.ServerConfigList = append(this.ServerConfigList, itemCfg)
		}
	}

	// 取得客户端配置
	{
		configNacosClient := client
		this.ClientConfig.TimeoutMs = gjson.Get(configNacosClient, "TimeoutMs").Uint()
		this.ClientConfig.ListenInterval = gjson.Get(configNacosClient, "ListenInterval").Uint()
		this.ClientConfig.BeatInterval = gjson.Get(configNacosClient, "BeatInterval").Int()
		this.ClientConfig.NamespaceId = gjson.Get(configNacosClient, "NamespaceId").String()
		this.ClientConfig.Endpoint = gjson.Get(configNacosClient, "Endpoint").String()

		this.ClientConfig.RegionId = gjson.Get(configNacosClient, "RegionId").String()
		this.ClientConfig.AccessKey = gjson.Get(configNacosClient, "AccessKey").String()
		this.ClientConfig.SecretKey = gjson.Get(configNacosClient, "SecretKey").String()
		this.ClientConfig.OpenKMS = gjson.Get(configNacosClient, "OpenKMS").Bool()

		this.ClientConfig.CacheDir = gjson.Get(configNacosClient, "CacheDir").String()
		this.ClientConfig.UpdateThreadNum, _ = strconv.Atoi(strconv.FormatInt(gjson.Get(configNacosClient, "UpdateThreadNum").Int(), 10))
		this.ClientConfig.NotLoadCacheAtStart = gjson.Get(configNacosClient, "NotLoadCacheAtStart").Bool()
		this.ClientConfig.UpdateCacheWhenEmpty = gjson.Get(configNacosClient, "UpdateCacheWhenEmpty").Bool()
		this.ClientConfig.Username = gjson.Get(configNacosClient, "Username").String()
		this.ClientConfig.Password = gjson.Get(configNacosClient, "Password").String()

		this.ClientConfig.LogDir = gjson.Get(configNacosClient, "LogDir").String()
		this.ClientConfig.RotateTime = gjson.Get(configNacosClient, "RotateTime").String()
		this.ClientConfig.MaxAge = gjson.Get(configNacosClient, "MaxAge").Int()
		this.ClientConfig.LogLevel = gjson.Get(configNacosClient, "LogLevel").String()
	}

	{
		//创建配置客户端..
		configClient, err := clients.CreateConfigClient(map[string]interface{}{
			"serverConfigs": this.ServerConfigList,
			"clientConfig":  this.ClientConfig,
		})
		if err != nil {
			return err
		}
		this.configClient = configClient
	}

	{
		//创建名字客户端..
		namingClient, err := clients.CreateNamingClient(map[string]interface{}{
			"serverConfigs": this.ServerConfigList,
			"clientConfig":  this.ClientConfig,
		})
		if err != nil {
			return err
		}
		this.namingClient = namingClient
	}

	return nil
}

func (this *DiscoveryContent) Register(serverName string, ip string, port uint64, weight float64) (bool, error) {

	if weight <= 0.0001 {
		weight = 10.01
	}

	success, err := this.namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serverName,
		Weight:      weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	if err != nil {
		return false, err
	}

	return success, nil
}

func (this *DiscoveryContent) Deregister(server *Service) error {

	return nil
}

func (this *DiscoveryContent) GetServers(name string) ([]*Service, error) {

	return nil, nil
}

func (this *DiscoveryContent) GetServerOne(serviceName string) (*Service, error) {
	srv := &Service{}
	ins, err := this.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})

	srv.Ins = ins

	return srv, err
}

func (this *DiscoveryContent) ListServices() ([]*Service, error) {

	return nil, nil
}

func (this *DiscoveryContent) String() string {

	return ""
}

func (this *DiscoveryContent) ListenConfig(dataId string, group string, callBack ListenerOption) string {

	return ""
}
