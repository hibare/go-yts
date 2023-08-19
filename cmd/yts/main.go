package main

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/history"
	"github.com/hibare/go-yts/internal/notifiers"
	log "github.com/sirupsen/logrus"
)

func ticker() {
	log.Info("[Start] Scraper task")

	movies := history.Movies{}
	urls := []string{"https://yts.mx/", "https://yts.autos/", "https://yts.rs/", "https://yts.lt/", "https://yts.do/"}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
	)

	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.Current.HTTPConfig.RequestTimeout,
			DualStack: true,
		}).DialContext,
	})

	c.SetRequestTimeout(config.Current.HTTPConfig.RequestTimeout)

	c.OnHTML("#popular-downloads", func(e *colly.HTMLElement) {
		temp := history.Movie{}
		e.ForEach("div .browse-movie-wrap", func(_ int, el *colly.HTMLElement) {
			temp.Link = el.ChildAttr(".browse-movie-link", "href")
			temp.TimeStamp = time.Now()
			temp.Title = el.ChildText(".browse-movie-title")
			temp.Year = el.ChildText(".browse-movie-year")
			temp.CoverImage = el.ChildAttr("img", "src")

			base, err := url.Parse(temp.Link)
			if err != nil {
				log.Fatal(err)
			}

			ref, err := url.Parse(temp.CoverImage)
			if err != nil {
				log.Fatal(err)
			}

			u := base.ResolveReference(ref)
			temp.CoverImage = u.String()

			movies[temp.Title] = temp
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referrer", "https://www.google.com/")
		log.Infof("Visiting URL: %s", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		log.Infof("Finished URL: %s", r.Request.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Infof("visited URL: %s Status Code: %d", r.Request.URL.String(), r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Errorf("Failed to load URL: %s Error: %s", r.Request.URL.String(), err)
	})

	q, _ := queue.New(
		2, &queue.InMemoryQueueStorage{MaxSize: 100},
	)

	for _, url := range urls {
		q.AddURL(url)
	}

	q.Run(c)

	log.Infof("Scraped %d movies", len(movies))

	h := history.ReadHistory(config.Current.StorageConfig.DataDir, config.Current.StorageConfig.HistoryFile)
	diff := history.DiffHistory(movies, h)
	history.WriteHistory(diff, h, config.Current.StorageConfig.DataDir, config.Current.StorageConfig.HistoryFile)
	log.Infof("Found %d new movies", len(diff))

	notifiers.Notify(diff)

	log.Info("[End] Scraper task")
}

func main() {
	initLogger()
	config.LoadConfig()
	log.Infof("Cron %s", config.Current.Schedule)
	log.Infof("Request Timeout %v", config.Current.HTTPConfig.RequestTimeout)
	log.Infof("Data directory %s", config.Current.StorageConfig.DataDir)
	log.Infof("History file %s", config.Current.StorageConfig.HistoryFile)
	log.Info("Starting scheduler")

	s := gocron.NewScheduler(time.UTC)
	s.Cron(config.Current.Schedule).Do(ticker)
	s.StartBlocking()
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
