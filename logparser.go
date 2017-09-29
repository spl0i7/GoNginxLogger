package main

import (
	"regexp"
	"github.com/oschwald/geoip2-golang"
	"github.com/tobie/ua-parser/go/uaparser"
	"time"
	"net"
	"strconv"
	"errors"
)

type Location struct {
	Latitude  float64
	Longitude float64
}
type IPInfo struct {
	Continent   string
	IP          string
	Country     string
	CountryCode string
	City        string
	Location    Location
}
type Record struct {
	RemoteIP       IPInfo
	Timestamp      time.Time
	RequestingUser string
	Method         string
	Request        string
	HttpVersion    string
	ResponseCode   string
	Size           int
	Referrer       string
	Client         *uaparser.Client
}

// [16/Apr/2017:11:37:06 -0400]
const LAYOUT = "02/Jan/2006:15:04:04 -0700"
const GEO_DB = "GeoLite2-City.mmdb"
const REGEX_DB = "regexes.yaml"

var geoDB *geoip2.Reader
var parser *uaparser.Parser

func initParser() {
	parser, _ = uaparser.New(REGEX_DB)
	db, err := geoip2.Open(GEO_DB)
	geoDB = db
	if err != nil {
		panic(err)
	}
}
func getTime(timestamp string) time.Time{
	res, err := time.Parse(LAYOUT, timestamp)
	if err != nil {
		panic("Invalid Time")
	}
	return res
}
func toInt(str string) int{
	res, err := strconv.Atoi(str)
	if err != nil {
		panic("Failed to convert to Integer")
	}
	return res
}
func parseIP(addr string) IPInfo{
	ip := net.ParseIP(addr)
	record, err := geoDB.City(ip)
	if err != nil {
		panic(err)
	}
	return IPInfo{
		City:        record.City.Names["en"],
		Country:     record.Country.Names["en"],
		Continent:   record.Continent.Names["en"],
		CountryCode: record.Country.IsoCode,
		IP:          addr,
		Location: Location{
			Longitude: record.Location.Longitude,
			Latitude:  record.Location.Latitude,
		},
	}

}
func getDocument(line string) (Record, error){

	regex := regexp.MustCompile(`(?P<remote_ip>\S*)\s-\s(?P<requesting_user>\S*)\s\[(?P<Timestamp>.*?)\]\s\"(?P<Method>\S*)\s*(?P<Request>\S*)\s*(HTTP\/)*(?P<http_version>.*?)\"\s(?P<response_code>\d{3})\s(?P<Size>\S*)\s\"(?P<Referrer>[^\"]*)\"\s\"(?P<Client>[^\"]*)`)
	match := regex.FindStringSubmatch(line)
	if len(match) == 12 {
		return Record{
			RemoteIP:       parseIP(match[1]),
			RequestingUser: match[2],
			Timestamp:      getTime(match[3]),
			Method:         match[4],
			Request:        match[5],
			HttpVersion:    match[7],
			ResponseCode:   match[8],
			Size:           toInt(match[9]),
			Referrer:       match[10],
			Client:         parser.Parse(match[11]),
		}, nil
	}else {
		return Record{}, errors.New("regex failed")
	}
}
