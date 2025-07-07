package rss

import (
	"encoding/xml"
	"github.com/grokify/html-strip-tags-go"
	"io/ioutil"
	"net/http"
	"news-aggregator/pkg/storage"
	"strings"
	"time"
)

type Xml struct {
	XMLName xml.Name `xml:"rss"`
	Chanel  Channel  `xml:"channel"`
}
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Posts       []Post `xml:"item"`
}
type Post struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

func Parse(handler func(url string) (resp *http.Response, err error), url string) ([]storage.Post, error) {
	resp, err := handler(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var xmlRes Xml
	err = xml.Unmarshal(body, &xmlRes)
	if err != nil {
		return nil, err
	}

	var data []storage.Post
	for _, item := range xmlRes.Chanel.Posts {
		var p storage.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.Link = item.Link
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
			if err != nil {
				return nil, err
			}
		}
		p.PubTime = t.Unix()
		data = append(data, p)
	}
	return data, nil
}
