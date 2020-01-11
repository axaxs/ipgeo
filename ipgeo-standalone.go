package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const BEHIND_PROXY = true
const IP_HEADER = "X-Real-IP"

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:9555", nil)
}

func ipfwd(r *http.Request) (string, []string) {
	remIP := strings.Split(r.RemoteAddr, ":")[0]
	remfwd := strings.Split(r.Header.Get("X-Forwarded-For"), ", ")
	if BEHIND_PROXY {
		remIP = r.Header.Get(IP_HEADER)
	}

	return remIP, remfwd
}

func extip(r *http.Request) string {
	ip, fwds := ipfwd(r)
	answer := ip
	if len(fwds) > 0 && fwds[0] != "" {
		answer = fmt.Sprintf("%s fwded for %s", ip, strings.Join(fwds, ", "))
	}

	return answer + "\n"
}

func geo(r *http.Request) string {
	ip, fwds := ipfwd(r)
	if len(fwds) > 0 {
		ip = fwds[len(fwds)-1]
	}

	if len(r.URL.Query()["ip"]) != 0 {
		ip = r.URL.Query()["ip"][0]
	}

	return get_coords(ip)
}

func get_coords(ip string) string {
	city, err := lookupIP(ip)
	if err != nil {
		return "\n"
	}

	ns := "N"
	ew := "E"
	if city.Location.Latitude < 0 {
		ns = "S"
		city.Location.Latitude *= -1
	}

	if city.Location.Longitude < 0 {
		ew = "W"
		city.Location.Longitude *= -1
	}

	return fmt.Sprintf("%.4f %s, %.4f %s\n", city.Location.Latitude, ns, city.Location.Longitude, ew)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var resp string
	switch strings.ToLower(strings.Split(r.Host, ":")[0]) {
	case "ip.lx.lc":
		resp = extip(r)
	case "geo.lx.lc", "geolocation.lx.lc":
		resp = geo(r)
	default:
		resp = "Unknown Request"
	}

	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Length", strconv.Itoa(len(resp)))
	w.Header().Add("Connection", "close")
	io.WriteString(w, resp)
}
