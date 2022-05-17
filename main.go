package main

import (
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/hedwi/endless"

	//"github.com/gin-contrib/pprof"
	// "github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
	"restdoc/config"
	redispool "restdoc/internal/database/redis"
	"restdoc/internal/database/snowflake"
	"restdoc/internal/middlewares"
	"restdoc/internal/route"
	"restdoc/logger"
)

//go:embed static
//go:embed wellknown
//go:embed restdoc
//go:embed templates
var static embed.FS

func init() {

	/*
		        if err := sentry.Init(sentry.ClientOptions{
		            Dsn: "https://022afd69e69b48f38eb6ecdf1ba44bf2@sty.hedwi.com/3",
		        }); err != nil {
		            fmt.Printf("Sentry initialization failed: %v\n", err)
		        }

				//defer sentry.Flush(2 * time.Second)
	*/
}

func loadCerts(dir string) (*tls.Config, error) {

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{}}
	for _, info := range mustReadDir(dir) {
		name := info.Name()
		dir := filepath.Join(dir, name)
		if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
			// Skip non-directories.
			continue
		}

		//certPath := filepath.Join(dir, "fullchain.pem")
		certPath := filepath.Join(dir, "cert.crt")
		glog.Info(certPath)
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			continue
		}
		keyPath := filepath.Join(dir, "cert.key")
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			continue
		}

		cer, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			glog.Error(err)
			//os.Exit(2)
			continue
		}

		//	s.tlsConfig.Certificates = append(s.tlsConfig.Certificates, cert)
		tlsConfig.Certificates = append(tlsConfig.Certificates, cer)
	}
	return tlsConfig, nil
}

func mustReadDir(path string) []os.FileInfo {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading %q directory: %v", path, err)
	}
	if len(dirs) == 0 {
		log.Fatalf("No entries found in %q", path)
	}

	return dirs
}

func main() {

	//if err := config.InitConfigInfo(); err != nil {
	if err := config.InitWithEnv(); err != nil {
		fmt.Println(err)
		panic("init config fail")
	}
	err := config.CheckConfig()
	if err != nil {
		panic(err)
	}

	config.DefaultConfig.VersionTimestamp = fmt.Sprintf("%d", time.Now().UnixNano())

	if config.DefaultConfig.Debug == true {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	isSaaS := false
	if config.DefaultConfig.SaaSDomain != "" {
		isSaaS = true
	}

	modelConfig := Models.ModelConfig{
		Debug:      config.DefaultConfig.Debug,
		IsSaaS:     isSaaS,
		Mysql:      config.DefaultConfig.Mysql.Host,
		Postgresql: config.DefaultConfig.Postgresql.Host,
	}

	Models.Init(&modelConfig)
	Models.CreateTables()

	if !isSaaS { //selfhost
		members := Models.MembersCount()
		if members > int64(config.DefaultConfig.Count) {
			err := errors.New("users limit error")
			panic(err)
		}
	}

	path := "logs"
	mode := os.ModePerm
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, mode)
	}
	file, _ := os.Create(strings.Join([]string{path, "log.txt"}, "/"))
	defer file.Close()
	loger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetDefault(loger)

	redispool.Init()
	snowflake.Init()

	/*
		err = utils.InitBucket()
		if err != nil {
			glog.Error(err)
			panic(err)
			return
		}
	*/

	//seelog.See(strings.Join([]string{path, "info.log"}, "/"),23456)

	middlewares.StaticBox = static
	middlewares.Init()

	route.StaticBox = static
	route.TemplateBox = static

	r := route.InitRouter()

	addr := config.DefaultConfig.Addr
	isTls := false
	if strings.HasSuffix(addr, ":443") {
		isTls = true
	}

	certDir := config.DefaultConfig.CertDir

	fmt.Println(config.DefaultConfig.Debug)
	if config.DefaultConfig.Debug {
		/*
			pprof.Register(r, &pprof.Options{
				// default is "debug/pprof"
				RoutePrefix: "logs/pprof",
			})
		*/
		s := &http.Server{
			Addr:           addr,
			Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		if isTls {
			tlsConfig, err := loadCerts(certDir)

			if err != nil {
				fmt.Println(err)
			} else {
				s.TLSConfig = tlsConfig
			}

			s.ListenAndServeTLS("", "") // listen and serve on 0.0.0.0:8080
		} else {
			s.ListenAndServe() // listen and serve on 0.0.0.0:8080
		}
	} else {
		server := endless.NewServer(addr, r)
		server.BeforeBegin = func(add string) {
			fmt.Printf("actual pid is %d\n", syscall.Getpid())
		}
		if isTls {

			tlsConfig, err := loadCerts(certDir)
			if err != nil {
				glog.Error(err)
			} else {
				server.TLSConfig = tlsConfig
			}
			err = server.ListenAndServeTLS("", "")
			if err != nil {
				glog.Error(err)
			}
		} else {

			err = server.ListenAndServe()
			if err != nil {
				glog.Error(err)
			}
		}

	}
}
