package main

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var mmr *geoip2.Reader
var sm sync.Mutex

func init() {
	f, err := os.Open("/srv/apps/GeoLite2-City.mmdb.gz")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer gzr.Close()
	data, err := ioutil.ReadAll(gzr)
	if err != nil {
		log.Fatalf(err.Error())
	}
	mmr, err = geoip2.FromBytes(data)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func lookupIP(ip string) (*geoip2.City, error) {
	sm.Lock()
	defer sm.Unlock()
	d := net.ParseIP(ip)
	return mmr.City(d)
}
