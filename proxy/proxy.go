package main

import (
	"github.com/armon/go-socks5"
	"log"
	"github.com/dimaxgl/ssproxy/store"
	"flag"
	"github.com/pkg/errors"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/dimaxgl/ssproxy/proxy/config"
	_ "github.com/dimaxgl/ssproxy/store/database"
)

var (
	confPath = flag.String(`c`, ``, `config file path`)
)

func main() {
	flag.Parse()

	storeConf, err := loadConfig(*confPath)
	if err != nil {
		log.Fatalln(`failed to load config`, err)
	}

	conf := &socks5.Config{}

	storeInst, err := store.Open(storeConf.Store.Type, storeConf.Store.Params)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(storeInst)

	//auth, err := store.NewDBStore(storeConf.DriverName, storeConf.Dsn)
	//if err != nil {
	//	log.Fatalln(`failed to create db store`, err)
	//}

	initUsers(storeInst)

	conf.Credentials = storeInst
	s, err := socks5.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	if err = s.ListenAndServe(`tcp`, storeConf.ListenAddress); err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(configPath string) (*config.Config, error) {
	if configPath == `` {
		return nil, errors.New(`empty config path`)
	}

	confBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, `failed to open config`)
	}

	var c config.Config

	if err = yaml.Unmarshal(confBytes, &c); err != nil {
		return nil, errors.Wrap(err, `failed to read yaml config`)
	}
	return &c, nil
}

func initUsers(s store.Store) {
	// here you can add users on every run of your proxy
	log.Println(s.Add(`user`,`password`))
}
