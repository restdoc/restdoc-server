package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Debug            bool
	Timeout          int
	SaaSDomain       string
	SelfDomain       string
	Count            int
	Addr             string
	CertDir          string
	Mysql            SqlDB
	Postgresql       SqlDB
	DefaultExpire    int
	SessionRedis     RedisInfo
	CacheRedis       RedisInfo
	Sources          map[string]string
	SignFields       map[string][]string
	IgnoreSignConfig map[string]bool
	TldCacheFile     string
	VersionTimestamp string
	DNSServer        string
	Server           string
	APIBaseUrl       string
	APIKey           string
	FromUser         string
	FromPassword     string
	SmtpServer       string
	SmtpServerPort   int
	SmtpTls          bool
	ToUser           string
	SupportUser      string
	AliPay           AliPayInfo
	GravatarUrl      string
	AllowOrigin      string
	Consul           ConsulInfo
}

type KVInfo struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	BucketName   string
	UseSSL       bool
}

type SqlDB struct {
	Host string
	PORT int
	DB   string
}

type RedisInfo struct {
	Address  string
	Password string
	DB       int
	PoolSize int
}

type HTTPSever struct {
	Addr string
}

type ConsulInfo struct {
	Addr string
}

type AliPayInfo struct {
	AppId      string
	PublicKey  string
	PrivateKey string
}

var DefaultConfig Config

func init() {
	//InitConfigInfo()
	InitWithEnv()
}

func InitConfigInfo() error {
	var confFileName = "./config.toml"
	_, err := toml.DecodeFile(confFileName, &DefaultConfig)
	return err
}

func InitWithEnv() error {

	_debug := os.Getenv("RESTDOC_DEBUG")
	debug := _debug == "true"
	_timeout := os.Getenv("RESTDOC_TIMEOUT") //timeout            int
	timeout, err := strconv.Atoi(_timeout)
	if err != nil {
		timeout = int(5)
	}

	saasDomain := os.Getenv("RESTDOC_DOMAIN") //saas domain

	httpAddr := os.Getenv("RESTDOC_HTTP_ADDR")            //Addr
	certDir := os.Getenv("RESTDOC_CERT_DIR")              //CertDir
	mysqlAddr := os.Getenv("RESTDOC_MYSQL")               //Mysql
	postgresqlAddr := os.Getenv("RESTDOC_POSTGRESQL")     //postgresql
	consulAddr := os.Getenv("RESTDOC_CONSUL")             //Consul
	_defaultExpire := os.Getenv("RESTDOC_DEFAULT_EXPIRE") //DefaultExpire int
	defaultExpire, err := strconv.Atoi(_defaultExpire)
	if err != nil {
		defaultExpire = int(86400 * 7)
	}

	sessionHedisAddr := os.Getenv("RESTDOC_SESSION_HEDIS_HOST")         //session_redis host
	sessionHedisPassword := os.Getenv("RESTDOC_SESSION_HEDIS_PASSWORD") //session_redis password
	_sessionHedisDB := os.Getenv("RESTDOC_SESSION_HEDIS_DB")            //session_redis db
	sessionHedisDB, err := strconv.Atoi(_sessionHedisDB)
	if err != nil {
		sessionHedisDB = 0
	}
	_sessionHedisPoolSize := os.Getenv("RESTDOC_SESSION_HEDIS_POOLSIZE") //session_redis poolsize
	sessionHedisPoolSize, err := strconv.Atoi(_sessionHedisPoolSize)
	if err != nil {
		sessionHedisPoolSize = 10
	}

	cacheHedisAddr := os.Getenv("RESTDOC_CACHE_HEDIS_HOST")         //cache_redis host
	cacheHedisPassword := os.Getenv("RESTDOC_CACHE_HEDIS_PASSWORD") //cache_redis password
	_cacheHedisDB := os.Getenv("RESTDOC_CACHE_HEDIS_DB")            //cache_redis db
	cacheHedisDB, err := strconv.Atoi(_cacheHedisDB)
	if err != nil {
		cacheHedisDB = 0
	}
	_cacheHedisPoolSize := os.Getenv("RESTDOC_CACHE_HEDIS_POOLSIZE") //cache_redis poolsize
	cacheHedisPoolSize, err := strconv.Atoi(_cacheHedisPoolSize)
	if err != nil {
		cacheHedisPoolSize = 10
	}

	//tldCacheFile := os.Getenv("RESTDOC_TLD_CACHE_FILE") //TldCacheFile

	dnsServer := os.Getenv("RESTDOC_DNSSERVER")
	if dnsServer == "" {
		dnsServer = "119.29.29.29:53"
	}
	server := os.Getenv("RESTDOC_SERVER")
	apiBaseUrl := os.Getenv("RESTDOC_APIBASEURL")
	apiKey := os.Getenv("RESTDOC_APIKEY")
	fromUser := os.Getenv("RESTDOC_FROMUSER")
	toUser := os.Getenv("RESTDOC_TOUSER")
	supportUser := os.Getenv("RESTDOC_SUPPORTUSER")
	allowOrigin := os.Getenv("RESTDOC_ALLOW_ORIGIN")

	alipayAppID := os.Getenv("RESTDOC_ALIPAY_APPID")
	alipayPublicKey := os.Getenv("RESTDOC_ALIPAY_PUBLIC_KEY")
	alipayPrivateKey := os.Getenv("RESTDOC_ALIPAY_PRIVATE_KEY")
	DefaultConfig = Config{
		Debug:         debug,
		Timeout:       timeout,
		SaaSDomain:    saasDomain,
		Addr:          httpAddr,
		CertDir:       certDir,
		Mysql:         SqlDB{Host: mysqlAddr},
		Postgresql:    SqlDB{Host: postgresqlAddr},
		DefaultExpire: defaultExpire,
		SessionRedis:  RedisInfo{Address: sessionHedisAddr, Password: sessionHedisPassword, DB: sessionHedisDB, PoolSize: sessionHedisPoolSize},
		CacheRedis:    RedisInfo{Address: cacheHedisAddr, Password: cacheHedisPassword, DB: cacheHedisDB, PoolSize: cacheHedisPoolSize},
		TldCacheFile:  "",
		DNSServer:     dnsServer,
		Server:        server,
		APIBaseUrl:    apiBaseUrl,
		APIKey:        apiKey,
		FromUser:      fromUser,
		ToUser:        toUser,
		SupportUser:   supportUser,
		AliPay:        AliPayInfo{AppId: alipayAppID, PublicKey: alipayPublicKey, PrivateKey: alipayPrivateKey},
		AllowOrigin:   allowOrigin,
		Consul:        ConsulInfo{Addr: consulAddr},
	}
	return nil
}

func CheckConfig() error {
	var err error
	if DefaultConfig.Addr == "" {
		err = errors.New("RESTDOC_HTTP_ADDR: Invalid addr ")
		return err
	}

	if DefaultConfig.Mysql.Host == "" && DefaultConfig.Postgresql.Host == "" {
		err = errors.New("RESTDOC_MYSQL and RESTDOC_POSTGRESQL are empty")
		return err
	}

	if DefaultConfig.SessionRedis.Address == "" {
		err = errors.New("RESTDOC_SESSION_HEDIS_HOST: Invalid sessionHedis address")
		return err
	}

	if DefaultConfig.SessionRedis.Password == "" {
		//err = errors.New("RESTDOC_SESSION_HEDIS_PASSWORD: Invalid sessionHedis password")
	}

	if DefaultConfig.CacheRedis.Address == "" {
		err = errors.New("RESTDOC_CACHE_HEDIS_HOST: Invalid cachehedis address")
		return err
	}

	if DefaultConfig.CacheRedis.Password == "" {
		//err = errors.New("RESTDOC_CACHE_HEDIS_PASSWORD: Invalid cachehedis address")
		//return err
	}

	if DefaultConfig.DNSServer == "" {
	}

	if DefaultConfig.Server == "" {
	}

	if DefaultConfig.APIBaseUrl == "" {
	}

	if DefaultConfig.APIKey == "" {
	}

	if DefaultConfig.FromUser == "" {
	}

	if DefaultConfig.ToUser == "" {
	}

	if DefaultConfig.SupportUser == "" {
	}

	if DefaultConfig.AliPay.AppId == "" {
	}

	if DefaultConfig.AliPay.PublicKey == "" {
	}

	if DefaultConfig.AliPay.PrivateKey == "" {
	}

	if DefaultConfig.GravatarUrl == "" {
	}

	if DefaultConfig.AllowOrigin == "" {
	}

	return nil
}
