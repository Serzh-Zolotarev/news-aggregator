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
	postsChan := make(chan []storage.NewsShortDetailed)
	errChan := make(chan error)

	for _, url := range conf.URLS {
		go func() {
			for {
				news, err := rss.Parse(http.Get, url)
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

	log.Println("[*] HTTP server is started on :8080")
	err = http.ListenAndServe(":8080", apiDb.Router())
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[*] HTTP server has been stopped. Reason: got sigterm")
	}
}
