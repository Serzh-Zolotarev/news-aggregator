package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"news-aggregator/pkg/api"
	"news-aggregator/pkg/rss"
	"news-aggregator/pkg/storage"
	"news-aggregator/pkg/storage/postgres"
	"time"
)

type server struct { // TODO
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

	// каналы для новостей и ошибок
	postsChan := make(chan []storage.Post)
	errChan := make(chan error)

	for _, url := range conf.URLS {
		go func() {
			for {
				news, err := rss.Parse(url)
				if err != nil {
					errChan <- err
					continue
				}
				postsChan <- news
				time.Sleep(time.Minute * time.Duration(conf.Period))
			}
		}()
	}

	go func() {
		for posts := range postsChan {
			err := db.AddPosts(posts)
			if err != nil {
				errChan <- err
				continue
			}
		}
	}()

	go func() {
		for err := range errChan {
			log.Println(err)
		}
	}()

	apiDb := api.New(db)

	err = http.ListenAndServe(":80", apiDb.Router())
	if err != nil {
		log.Fatal(err)
	}
}
