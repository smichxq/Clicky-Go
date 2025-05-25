module clicky.website/clicky/security

go 1.24.2

replace clicky.wesite/clicky/common/serversuite => ../../common/

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require github.com/redis/go-redis/v9 v9.8.0

require (
	github.com/cloudwego/kitex v0.11.3 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	go.opentelemetry.io/otel v1.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/gorm v1.25.10 // indirect
)

require (
	github.com/apache/thrift v0.13.0 // indirect
	github.com/bytedance/gopkg v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/gopkg v0.1.3-0.20241115063537-a218fe69d609 // indirect
	github.com/cloudwego/kitex/pkg/protocol/bthrift v0.0.0-20250515033522-7c4ae57b7288 // indirect
	github.com/cloudwego/thriftgo v0.3.18 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/kitex-contrib/obs-opentelemetry/logging/logrus v0.0.0-20241120035129-55da83caab1b
)
