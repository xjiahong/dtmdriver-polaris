# dtmdriver-polaris


#### Background 背景
Polaris is a new product that Tencent Cloud developed for load balance and router for K8S etcs.

北极星是腾讯云目前推出的用于容器寻址和负载均衡的产品

### 20230810调整
升级了最新的北极星版本，且使用环境变量来获取北极星的注册中心配置如IP、端口

支持的环境变量：

| 环境变量                                     | 备注                      |
|:-----------------------------------------|:------------------------|
| DTM-MICRO-SERVICE-POLARIS-NAMESPACE      | 北极星配置-命名空间              |
| DTM-MICRO-SERVICE-POLARIS-HOST           | 北极星配置-服务监听host          |
| DTM-MICRO-SERVICE-POLARIS-PROVIDERPORT   | 北极星配置-服务实例监听port-服务注册   |
| DTM-MICRO-SERVICE-POLARIS-SERVICETOKEN   | 北极星配置-服务访问令牌            |



