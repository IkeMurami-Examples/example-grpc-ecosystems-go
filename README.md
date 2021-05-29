# grpc-ecosystems-example
Пример микросервиса с использованием grpc-ecosystems

# Поэтапно

## Init

Для создания CLI и поддержки конфигов
```
go get -u github.com/spf13/cobra
go get -u github.com/spf13/viper
```

Для умного логирования
```
go get -u github.com/rs/zerolog
go get -u github.com/prometheus/client_golang/prometheus/promhttp
go get -u github.com/getsentry/sentry-go
```

Для работы с protobuf-структурами
```
go get -u github.com/jhump/protoreflect
```

Работа с gRPC:
```
go get -u google.golang.org/grpc
go get -u github.com/grpc-ecosystem/grpc-gateway
```

Другое:
```
go get -u golang.org/x/sync
go get -u github.com/slok/go-http-metrics
go get -u github.com/hako/durafmt 
```

Трейсинг
```
go get -u github.com/uber/jaeger-client-go
go get -u github.com/uber/jaeger-client-go/config
go get -u github.com/opentracing-contrib/go-stdlib
```

## 
