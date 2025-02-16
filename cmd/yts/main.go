package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/google/uuid"
	commonContext "github.com/hibare/GoCommon/v2/pkg/context"
	commonErrors "github.com/hibare/GoCommon/v2/pkg/errors"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/db"
	"github.com/hibare/go-yts/internal/notifiers"
)

func ConstructURL(baseUrl *url.URL, refUrl string) (string, error) {
	ref, err := url.Parse(refUrl)
	if err != nil {
		return "", err
	}

	u := baseUrl.ResolveReference(ref)

	return u.String(), nil
}

func ticker(ctx context.Context) {
	ctx = context.WithValue(ctx, commonContext.ContextKey("run_id"), uuid.NewString())

	slog.InfoContext(ctx, "[Start] Scraper task")

	movies := map[string]db.Movies{}
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
		temp := db.Movies{}
		e.ForEach("div .browse-movie-wrap", func(_ int, el *colly.HTMLElement) {
			var err error
			temp.Link = el.ChildAttr(".browse-movie-link", "href")
			temp.Title = el.ChildText(".browse-movie-title")
			year, err := strconv.Atoi(el.ChildText(".browse-movie-year"))
			if err != nil {
				slog.WarnContext(ctx, "Failed to convert year", "error", err)
			}
			temp.Year = year
			temp.CoverImage = el.ChildAttr("img", "src")

			temp.Link, err = ConstructURL(e.Request.URL, temp.Link)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
				return
			}

			temp.CoverImage, err = ConstructURL(e.Request.URL, temp.CoverImage)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to construct URL", "error", err)
				return
			}

			movies[temp.Title] = temp
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referrer", "https://www.google.com/")
		slog.InfoContext(ctx, "Visiting URL", "url", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		slog.InfoContext(ctx, "Finished URL", "url", r.Request.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		slog.InfoContext(ctx, "Finished URL", "url", r.Request.URL.String(), "status_code", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		slog.ErrorContext(ctx, "Failed to load URL", "url", r.Request.URL.String(), "error", err)
	})

	q, _ := queue.New(
		2, &queue.InMemoryQueueStorage{MaxSize: 100},
	)

	for _, url := range urls {
		q.AddURL(url)
	}

	q.Run(c)

	slog.InfoContext(ctx, "Scraped movies", "total", len(movies))

	newMovies := []db.Movies{}

	for _, v := range movies {
		// Check if movie already exists
		slog.DebugContext(ctx, "processing movie", "movie", v.Title)
		_, err := db.GetMovieByTitle(context.Background(), v.Title)
		if err != nil {
			if errors.Is(err, commonErrors.ErrRecordNotFound) {
				slog.DebugContext(ctx, "New movie found", "movie", v.Title)
				// Its a new movie
				newMovies = append(newMovies, v)
				if err = db.AddMovie(context.Background(), v); err != nil {
					slog.ErrorContext(ctx, "Failed to add movie", "movie", v.Title, "error", err)
				}
				slog.DebugContext(ctx, "New movie added to DB", "movie", v.Title)
			} else {
				slog.ErrorContext(ctx, "Failed to get movie", "movie", v.Title, "error", err)
				continue
			}
		} else {
			slog.DebugContext(ctx, "existing movie found", "movie", v.Title)
			// Update last found & occurrence
			if err = db.UpdateMovieLastFound(context.Background(), v.Title); err != nil {
				slog.ErrorContext(ctx, "Failed to update last found", "movie", v.Title, "error", err)
			}

			slog.DebugContext(ctx, "Movie last found updated", "movie", v.Title)
		}
	}

	slog.InfoContext(ctx, "Found new movies", "total", len(newMovies))
	notifiers.Notify(ctx, newMovies)

	slog.InfoContext(ctx, "[End] Scraper task")
}

func main() {
	ctx := context.Background()
	commonLogger.InitDefaultLogger()
	config.LoadConfig()
	slog.InfoContext(ctx, "Config", "cron", config.Current.Schedule, "request_timeout", config.Current.HTTPConfig.RequestTimeout, "data_dir", config.Current.StorageConfig.DataDir)

	slog.InfoContext(ctx, "Migrating DB")
	if err := db.Migrate(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to migrate DB", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Starting scheduler")

	s := gocron.NewScheduler(time.UTC)
	s.Cron(config.Current.Schedule).Do(ticker, ctx)
	s.StartBlocking()
}
