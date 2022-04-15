module restdoc

go 1.18

replace github.com/blugelabs/bluge => /Users/solos/dev/bluge // 绝对路径 或 相对路径 都可以

replace github.com/hedwi/douceur => /Users/solos/dev/douceur // 绝对路径 或 相对路径 都可以

replace restdoc-models => /Users/solos/dev/restdoc-models

require (
	github.com/BobuSumisu/aho-corasick v1.0.3
	github.com/BurntSushi/toml v0.3.1
	github.com/MauriceGit/skiplist v0.0.0-20191117202105-643e379adb62
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/RoaringBitmap/roaring v0.6.1
	github.com/Xuanwo/go-locale v1.0.0
	github.com/andy-kimball/arenaskl v0.0.0-20200617143215-f701008588b9
	github.com/asggo/spf v0.0.0-20200529014219-3e270ddb6136
	github.com/benbjohnson/css v0.0.0-20141214004234-6538f8623a6a
	github.com/blugelabs/bluge v0.0.0-00010101000000-000000000000
	github.com/deckarep/golang-set v1.7.1
	github.com/ericlagergren/decimal v0.0.0-20210307182354-5f8425a47c58
	github.com/getsentry/sentry-go v0.10.0
	github.com/gin-contrib/multitemplate v0.0.0-20210428235909-8a2f6dd269a0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.1
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/uuid v1.2.0
	github.com/hashicorp/consul/api v1.3.0
	github.com/hedwi/douceur v0.0.0-00010101000000-000000000000
	github.com/hedwi/endless v0.0.0-20210910090835-db1cd8d23952
	github.com/huandu/skiplist v1.1.0
	github.com/jbenet/go-is-domain v1.0.5
	github.com/joeguo/tldextract v0.0.0-20180214020933-b623e0574407
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/mbobakov/grpc-consul-resolver v1.4.4
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/miekg/dns v1.1.41
	github.com/minio/minio-go/v7 v7.0.13
	github.com/nanmu42/gzip v1.1.0
	github.com/sean-public/fast-skiplist v0.0.0-20200308194023-d7f7945b944e
	github.com/smartystreets/goconvey v1.6.4
	github.com/szuecs/gin-glog v1.1.1
	github.com/zhangyunhao116/skipset v0.9.0
	github.com/zserge/lorca v0.1.9
	golang.org/x/net v0.0.0-20210323141857-08027d57d8cf
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.36.0
	gorm.io/gorm v1.23.2
	restdoc-models v0.0.0-00010101000000-000000000000
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/adamzy/cedar-go v0.0.0-20170805034717-80a9c64b256d // indirect
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/armon/go-metrics v0.3.2 // indirect
	github.com/axiomhq/hyperloglog v0.0.0-20191112132149-a4c4c47bc57f // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/mmap-go v1.0.2 // indirect
	github.com/blevesearch/segment v0.9.0 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blugelabs/bluge_segment_api v0.1.0 // indirect
	github.com/blugelabs/ice v0.1.1 // indirect
	github.com/caio/go-tdigest v3.1.0+incompatible // indirect
	github.com/chris-ramon/douceur v0.2.0 // indirect
	github.com/cockroachdb/cockroach-go/v2 v2.2.8 // indirect
	github.com/couchbase/vellum v1.0.2 // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/form v3.1.4+incompatible // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190910122728-9d188e94fb99 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.1.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.8.5 // indirect
	github.com/huichen/sego v0.0.0-20180617034105-3f3c8a8cfacc // indirect
	github.com/issue9/assert v1.5.0 // indirect
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
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.11.3 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/minio/md5-simd v1.1.0 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/mschoch/smat v0.0.0-20160514031455-90eadee771ae // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/signalsciences/ac v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	github.com/zhangyunhao116/sbconv v0.2.1 // indirect
	github.com/zhangyunhao116/wyhash v0.3.2 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.3.2 // indirect
	gorm.io/driver/postgres v1.3.4 // indirect
)
