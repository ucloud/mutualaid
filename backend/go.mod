module github.com/ucloud/mutualaid/backend

go 1.16

require (
	github.com/btcsuite/btcutil v1.0.2
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/go-kratos/kratos/contrib/metrics/prometheus/v2 v2.0.0-20211213133101-2e045c3e42e1
	github.com/go-kratos/kratos/v2 v2.3.1
	github.com/go-kratos/swagger-api v1.0.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/google/wire v0.5.0
	github.com/howeyc/crc16 v0.0.0-20171223171357-2b2a61e366a6
	github.com/imroc/req/v3 v3.10.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/prometheus/client_golang v1.11.0
	github.com/satori/go.uuid v1.2.0
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/viper v1.10.0 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/ucloud/mutualaid/backend/tools/micro-mini v0.10.5
	github.com/ucloud/mutualaid/backend/tools/ulibgo v0.10.1
	github.com/wechatpay-apiv3/wechatpay-go v0.2.11
	go.uber.org/zap v1.19.1
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gorm.io/driver/mysql v1.3.3
	gorm.io/gen v0.2.44
	gorm.io/gorm v1.23.4
)

replace (
	github.com/ucloud/mutualaid/backend/tools/micro-mini => ./tools/micro-mini
	github.com/ucloud/mutualaid/backend/tools/ulibgo => ./tools/ulibgo
)
