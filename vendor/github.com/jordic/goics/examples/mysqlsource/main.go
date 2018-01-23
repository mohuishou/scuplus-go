package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jordic/goics"
)

var (
	configFile = flag.String("config", "config.json", "Path to config file")
	Db         *sqlx.DB
)

var config struct {
	Dsn           string `json:"dsn"`
	ServerAddress string
}

const version = "0.1"

func main() {

	flag.Parse()
	load_config()
	//fmt.Printf("%v", config)
	var err error
	Db, err = sqlx.Connect("mysql", config.Dsn)
	if err != nil {
		panic("Cant connect to database")
	}

	m := mux.NewRouter()
	m.Handle("/version", http.HandlerFunc(Version))
	m.Handle("/limpieza", http.HandlerFunc(LimpiezaHandler))

	http.Handle("/", m)
	log.Print("Server Started")
	log.Fatal(http.ListenAndServe(config.ServerAddress, nil))
}

func Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version %s", version)
	return
}

func LimpiezaHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("Calendar request")
	// Setup headers for the calendar
	w.Header().Set("Content-type", "text/calendar")
	w.Header().Set("charset", "utf-8")
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("filename", "calendar.ics")
	// Get the Collection models
	collection := GetReservas()
	// Encode it.
	goics.NewICalEncode(w).Encode(collection)
}

// Load and parse config json file
func load_config() {
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading config file %s", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error decoding config file %s", err)
	}
}
