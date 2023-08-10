package driver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dtm-labs/dtmdriver"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"

	_ "github.com/polarismesh/grpc-go-polaris"
)

const (
	Name = "dtm-driver-polaris"
)

var (
	provider   api.ProviderAPI
	consumer   api.ConsumerAPI
	PolarisTtl = 5 * time.Second // 心跳周期
)

type (
	polarisDriver struct{}
	dialOptions   struct {
		Namespace string `json:"Namespace"`
	}
)

func (p *polarisDriver) GetName() string {
	return Name
}

// RegisterGrpcResolver grpc-go-polaris的init已完成注册
func (p *polarisDriver) RegisterGrpcResolver() {}

func firstIp() net.IP {
	ias, _ := net.InterfaceAddrs()
	for _, address := range ias {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP
		}
	}
	return nil
}

// RegisterGrpcService 向北极星注册dtm server服务
// target polaris://ip:port/service?namespace=[Test,Pre-release,Production]
func (p *polarisDriver) RegisterGrpcService(target, token string) error {
	var err error
	// 主要考虑是dtmsvr的初始化，如果使用其他driver，polaris的consumer和provider不都应该初始化
	if provider == nil {
		config := GetPolarisConfiguration(GetPolarisConf())
		provider, err = api.NewProviderAPIByConfig(config)
		if err != nil {
			return err
		}
	}
	//// token为空不注册，用于托管服务的场景
	//if token == "" {
	//	return nil
	//}
	u, err := url.Parse(target)
	if err != nil {
		return err
	}
	// namespace
	ns := u.Query().Get("namespace")
	if ns == "" {
		return fmt.Errorf("namespace not found %s", target)
	}

	// service
	if len(u.Path) <= 1 {
		return fmt.Errorf("service not found %s", target)
	}
	service := u.Path[1:]

	ip := net.ParseIP(u.Hostname())
	if ip == nil || ip.IsUnspecified() {
		ip = firstIp()
	}
	if ip == nil {
		return fmt.Errorf("ip not found %s", target)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return fmt.Errorf("parse port failed %s, target :%s", err.Error(), target)
	}

	request := &api.InstanceRegisterRequest{}
	{
		request.Namespace = ns
		request.Service = service
		request.ServiceToken = token
		request.Host = ip.String()
		request.Port = port
		request.SetTTL(int(PolarisTtl))
	}

	rsp, err := provider.Register(request)
	if rsp == nil || err != nil {
		return fmt.Errorf("register instance failed, err: %w", err)
	}

	hbReq := &api.InstanceHeartbeatRequest{
		InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
			Namespace:    ns,
			Service:      service,
			Host:         ip.String(),
			Port:         port,
			ServiceToken: token,
			InstanceID:   rsp.InstanceID,
		},
	}

	// 心跳上报&关闭的反注册
	go func() {
		if err = provider.Heartbeat(hbReq); nil != err {
			fmt.Println("polaris heartbeat error", err)
		}
		quit := make(chan os.Signal)
		signal.Notify(
			quit,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGTERM,
		)
		ticker := time.NewTicker(PolarisTtl)
		for {
			select {
			case <-ticker.C:
				if err = provider.Heartbeat(hbReq); nil != err {
					fmt.Println("polaris heartbeat error", err)
				}
			case <-quit:
				provider.Deregister(&api.InstanceDeRegisterRequest{
					InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
						Namespace:    ns,
						Service:      service,
						Host:         ip.String(),
						Port:         port,
						ServiceToken: token,
						InstanceID:   rsp.InstanceID,
					},
				})
				ticker.Stop()
				provider.Destroy()
			}
		}
	}()

	return nil
}

// ParseServerMethod 面向github.com/polarismesh/grpc-go-polaris解析
// uri polaris://service/package.service/method?namespace=Test
func (p *polarisDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	if !strings.Contains(uri, "//") { // 处理无scheme的情况，如果您没有直连，可以不处理
		sep := strings.IndexByte(uri, '/')
		if sep == -1 {
			return "", "", fmt.Errorf("bad url: '%s'. no '/' found", uri)
		}
		return uri[:sep], uri[sep:], nil

	}
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", fmt.Errorf("parse url failed, err: %w", err)
	}
	opts := &dialOptions{}
	opts.Namespace = u.Query().Get("namespace")
	jsonStr, _ := json.Marshal(opts)

	server = u.Scheme + "://" + u.Host + "/?options=" + base64.URLEncoding.EncodeToString(jsonStr)
	method = u.Path
	return
}

func init() {
	dtmdriver.Register(&polarisDriver{})
}
