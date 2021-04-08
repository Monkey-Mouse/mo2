# ARCHITECTURE

Mo2总体为微服务架构，生产环境使用docker-compose组合多个容器而成


## UML
```nomnoml
[<actor>user browser] <-> http2[docker containers|
    [nginx]<:--:>http2[mo2 main] 
    [mo2 main] <:--:> grpc[mo2 notification]
    [mo2 main] <:--:> grpc[mo2 img log]
    [mo2 main] <:--:> http2 [mongodb clusters]
    [mo2 notification] <:--:> http2 [mongodb clusters]
    [mo2 main] <:--:> http2 [redis]
    [mo2 audit log] <:--:> http2 [mongodb clusters]
    [mo2 audit log] <:--:> grpc [mo2 main]
    [mo2 main] <:--:> http2 [mo2 search]
]
```

![arch](imgs/architecture.png)
