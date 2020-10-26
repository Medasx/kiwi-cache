package main

import (
	"html/template"
	"kiwi-entry-task/cache"
	"kiwi-entry-task/fetcher"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const sqLiteFilePath = "currencies.db"

func main() {
	locFetcher := fetcher.NewLocationsFetcher()

	sqLiteFetcher, err := fetcher.NewSQLiteFetcher(sqLiteFilePath)
	if err != nil {
		log.Fatal(err)
	}
	locationCache := cache.NewCache(locFetcher, time.Minute, func(err error) {
		log.Println(err)
	})
	currencyCache := cache.NewCache(sqLiteFetcher, time.Minute, func(err error) {
		log.Println(err)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			value   string
			cache   cache.Cacher
			idStr   string
			id      int
			isError bool
		)

		switch r.URL.Query().Get("cache") {
		case "loc":
			cache = locationCache
		case "currency":
			cache = currencyCache
		}

		if idStr = r.URL.Query().Get("id"); idStr != "" {
			id, err = strconv.Atoi(idStr)
			if err != nil {
				value = err.Error()
				isError = true
			}
		}

		if cache != nil && id != -1 {
			val, err := cache.Get(id)
			if err != nil {
				value = err.Error()
				isError = true
			} else {
				value = val
			}
		}

		lp := filepath.Join("index.tmpl")

		tmpl, _ := template.ParseFiles(lp)
		err := tmpl.ExecuteTemplate(w, "index.tmpl", map[string]interface{}{
			"Value":   value,
			"Id":      idStr,
			"IsError": isError,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	})

	log.Println("Listening on :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
