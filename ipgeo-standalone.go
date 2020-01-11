package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:9555", nil)
}

func extip(r *http.Request) string {
	remIP := strings.Split(r.RemoteAddr, ":")[0]
	remfwd := r.Header.Get("X-Forwarded-For")
	var answer string
	if remIP == remfwd || remfwd == "" {
		answer = remIP
	} else {
		answer = fmt.Sprintf("%s fwded for %s", remIP, remfwd)
	}
	return answer + "\n"
}

func geo(r *http.Request) string {
	remIP := strings.Split(r.RemoteAddr, ":")[0]
	remfwd := r.Header.Get("X-Forwarded-For")
	var ip string
	if remIP == remfwd || remfwd == "" {
		ip = remIP
	} else {
		ip = remfwd
	}
	if strings.Contains(ip, ",") {
		ip = strings.TrimSpace(strings.Split(ip, ",")[len(ip)-1])
	}
	return get_coords(ip)
}

func get_coords(ip string) string {
	city, err := lookupIP(ip)
	if err != nil {
		return "\n"
	}
	return fmt.Sprintf("%.4f %.4f\n", city.Location.Latitude, city.Location.Longitude)
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
