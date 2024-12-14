package rssParser

import (
	"errors"
	"fmt"

	"github.com/mmcdole/gofeed"
)

var (
	source string = "https://lenta.ru/rss"
	feed   *gofeed.Feed
	fp     = gofeed.NewParser()
)

var (
	errNoNews = errors.New("нет информации о новостях")
)

func read() error {
	var err error
	feed, err = fp.ParseURL(source)
	if err != nil {
		return err
	}
	return nil
}

func print() error {
	if len(feed.Items) == 0 {
		return errNoNews
	}

	for i, item := range feed.Items {
		fmt.Println(i+1, ")", item.Title)
	}

	return nil
}

func RssParser() {
	err := read()

	if err != nil {
		fmt.Println("Произошла ошибка при считывании rss - ленты:", err)
		return
	}

	err = print()
	if err != nil {
		fmt.Println(err)
	}
}
