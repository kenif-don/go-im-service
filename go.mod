module go-im-service

go 1.22

toolchain go1.22.0

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/aws/aws-sdk-go v1.53.12
	github.com/go-netty/go-netty v1.6.5
	github.com/google/uuid v1.6.0
	github.com/h2non/filetype v1.1.3
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.9.3
	go-nio-client-sdk v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.34.1
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.10
)

replace go-nio-client-sdk => github.com/kenif-don/go-nio-client-sdk v1.0.3

require (
	github.com/go-netty/go-netty-transport v1.7.10 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
