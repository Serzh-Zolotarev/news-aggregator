package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"news-aggregator/pkg/api"
	"news-aggregator/pkg/storage"
	"news-aggregator/pkg/storage/postgres"
)

type server struct {
	db  storage.Interface
	api *api.API
}

type config struct {
	URLS       []string `json:"rss"`
	Period     int      `json:"request_period"`
	PostgresDb string   `json:"postgres_db"`
}

func main() {
	confFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var conf config
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		log.Fatal(err)
	}

	// подключение к БД
	db, err := postgres.New(conf.PostgresDb)
	if err != nil {
		log.Fatal(err)
	}

	apiDb := api.New(db)

	err = http.ListenAndServe(":80", apiDb.Router())
	if err != nil {
		log.Fatal(err)
	}
}
