package main

import (
	"fmt"
	"github.com/nranchev/go-libGeoIP"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const DEFAULT_PORT = 8000

const heroku = true

const dbFile = "geoip.dat"

var gi *libgeo.GeoIP

func get_database() {
	cmd := exec.Command("wget", "http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Getting Database from maxmind..\n")
	err = cmd.Wait()

	cmd = exec.Command("gzip", "-d", "GeoLiteCity.dat.gz")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("UnGzipping..\n")
	err = cmd.Wait()

	cmd = exec.Command("mv", "GeoLiteCity.dat", dbFile)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Renaming..\n")
	err = cmd.Wait()

	cmd = exec.Command("rm", "GeoLiteCity.dat.gz")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleting Archive..\n")
	err = cmd.Wait()

}

func geoip_init() {
	var err error

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		fmt.Printf("GeoIP Database does not exist: %s\n", dbFile)
		get_database()
	}

	fmt.Printf("Initialising GeoIP database\n")

	// Load the database file, exit on failure
	gi, err = libgeo.Load(dbFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	var addr string
	if !heroku {
		split := strings.Split(req.RemoteAddr, ":")
		addr = split[0]
	} else {
		addr = req.Header.Get("X-forwarded-For")
	}
	if len(addr) > 0 {
		loc := gi.GetLocationByIP(addr)

		collector := ""

		if loc != nil {
			collector = collector + fmt.Sprintf("Country: %s (%s)\n", loc.CountryName, loc.CountryCode)
			collector = collector + fmt.Sprintf("City: %s\n", loc.City)
			collector = collector + fmt.Sprintf("Region: %s\n", loc.Region)
			collector = collector + fmt.Sprintf("Postal Code: %s\n", loc.PostalCode)
			collector = collector + fmt.Sprintf("Latitude: %f\n", loc.Latitude)
			collector = collector + fmt.Sprintf("Longitude: %f\n", loc.Longitude)
		}

		fmt.Printf("ip: " + addr + "\n" + collector) // Logging

	} else {
		fmt.Fprintf(w, "unknown\n")
	}
}

func main() {
	var port int = DEFAULT_PORT
	var err error
	if len(os.Args) > 1 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	geoip_init()
	fmt.Printf("geoip-echo listening on port: %d\n", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
