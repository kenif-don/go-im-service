module IM-Service

go 1.21.3

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/glebarez/sqlite v1.10.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
	im-sdk v0.0.0-00010101000000-000000000000
)

replace im-sdk => github.com/don764372409/im-sdk v0.1.6

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/go-netty/go-netty v1.6.4 // indirect
	github.com/go-netty/go-netty-transport v1.7.5 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.3.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.13.0 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)
