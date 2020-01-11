package main

import ("sync";"os";"log";"net";"io/ioutil";"github.com/oschwald/geoip2-golang";"compress/gzip")

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
	defer mmr.Close()
}

func lookupIP(ip string) (*geoip2.City, error) {
	sm.Lock()
	defer sm.Unlock()
	d := net.ParseIP(ip)
	return mmr.City(d)
}

