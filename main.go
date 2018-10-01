package main

import (
    "flag"
    "log"
    "net"
    "net/http"
    "encoding/json"

    "github.com/ipipdotnet/datx-go"
)

var addr   = flag.String("addr", ":80", "service listen address")
var dbfile = flag.String("dbfile", "./17monipdb.datx", "db file of 17monipdb.datx path")

var city *datx.City

func init() {
    flag.Parse()

    var err error
    if city, err = datx.NewCity(*dbfile); err != nil {
        log.Fatal(err)
    }
}

func Handler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    ip := query.Get("ip")

    if ip == "" {
        ip, _, _ = net.SplitHostPort(r.RemoteAddr)
    }

    // loc, err := city.Find(ip)
    loc, err := city.FindLocation(ip)
    if err != nil {
        w.WriteHeader(404)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
    json.NewEncoder(w).Encode(loc)

    go func() {
        log.Printf("[%s] %s, %s", ip, r.URL, loc.ToJSON())
        }()
}

func main() {
    http.HandleFunc("/", Handler)
    log.Fatal(http.ListenAndServe(*addr, nil))
}