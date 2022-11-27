module restdoc

go 1.17

replace github.com/blugelabs/bluge => /Users/solos/dev/bluge // 绝对路径 或 相对路径 都可以

replace github.com/hedwi/douceur => /Users/solos/dev/douceur // 绝对路径 或 相对路径 都可以

replace restdoc-models => /Users/solos/dev/restdoc-models

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Xuanwo/go-locale v1.0.0
	github.com/deckarep/golang-set v1.7.1
	github.com/ericlagergren/decimal v0.0.0-20210307182354-5f8425a47c58
	github.com/getsentry/sentry-go v0.10.0
	github.com/gin-contrib/multitemplate v0.0.0-20210428235909-8a2f6dd269a0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/google/uuid v1.2.0
	github.com/hedwi/endless v0.0.0-20210910090835-db1cd8d23952
	github.com/jbenet/go-is-domain v1.0.5
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/nanmu42/gzip v1.1.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/szuecs/gin-glog v1.1.1
	golang.org/x/net v0.0.0-20210323141857-08027d57d8cf // indirect
	golang.org/x/text v0.3.7
	gorm.io/gorm v1.23.2
	restdoc-models v0.0.0-00010101000000-000000000000
)

require (
	github.com/amalfra/etag v1.0.0
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/cockroachdb/cockroach-go/v2 v2.2.8 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190910122728-9d188e94fb99 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/jackc/pgx/v4 v4.15.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.11.3 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/signalsciences/ac v1.2.0 // indirect
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.3.2 // indirect
	gorm.io/driver/postgres v1.3.4 // indirect
)
