# go-tools
hy_encrypt：
    aes cbc 模式加解密
    MD5
    rsa 加解密
    tea 加解密
    
hy_file：获取文件信息相关函数

hy_gmsm：过密加密方式
    sm2：非对称加密
    sm3：类似md5
    sm4：对称加密 包含 cbc  cfb ecb  ofb
    
hy_http: http客户端封装，get和post  post包含json  form byte

hy_imgutils: 获取缩略图和图片信息

hy_log:
    logrus 简单封装
    zap 简单封装
    
hy_mq: 消息队列封装
    kafka 简单封装
    rocketmq 简单封装
    
hy_mysql:
    gorm: 单机和主从连接封装
    xorm: 单机和主从连接封装
    
hy_servicectrl: 微服务常用服务治理工具
    hy_hystrix: 熔断
    hy_prometheus：普罗米修斯
    hy_ratelimit：限流
    hy_tracer：链路追踪
    
hymongodb: mongodb客户端封装

hyredis：redigo的的封装
