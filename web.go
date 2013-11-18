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

func run_command(command string, msg string, args []string) {
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg)
	err = cmd.Wait()

}

func get_database() {
	run_command("wget", "Getting Database from maxmind..", []string{"http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz"})
	run_command("gzip", "UnGzipping..", []string{"-d", "GeoLiteCity.dat.gz"})
	run_command("mv", "Renaming..", []string{"GeoLiteCity.dat", dbFile})
	run_command("rm", "Deleting Archive..", []string{"GeoLiteCity.dat.gz"})
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

		resp := "ip: " + addr + "\n" + collector
		fmt.Fprintf(w, resp)
		fmt.Printf(resp)

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
