package db

import (
	"context"
	"embed"
	"errors"
	"path/filepath"
	"time"

	commonDB "github.com/hibare/GoCommon/v2/pkg/db"
	commonErrors "github.com/hibare/GoCommon/v2/pkg/errors"
	"github.com/hibare/go-yts/internal/config"
	"github.com/hibare/go-yts/internal/constants"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Movies struct {
	Title        string `gorm:"primaryKey"`
	Link         string
	CoverImage   string
	Year         int
	FirstFoundOn time.Time `gorm:"autoCreateTime"`
	LastFoundOn  time.Time
	Occurrences  int
}

func (Movies) TableName() string {
	return "movies"
}

var getDSN = func() string {
	return filepath.Join(config.Current.StorageConfig.DataDir, constants.DefaultSQLiteDB)
}

func getNewDBClient(ctx context.Context) (*commonDB.DB, error) {
	dbConfig := commonDB.DatabaseConfig{
		DSN:            getDSN(),
		MigrationsPath: "migrations",
		MigrationsFS:   migrationsFS,
		DBType:         &commonDB.SQLiteDatabase{},
	}
	return commonDB.NewClient(ctx, dbConfig)
}

func Migrate(ctx context.Context) error {
	client, err := getNewDBClient(ctx)
	if err != nil {
		return err
	}

	return client.Migrate()
}

func GetAllMovies(ctx context.Context) ([]Movies, error) {
	client, err := getNewDBClient(ctx)
	if err != nil {
		return nil, nil
	}

	var movies []Movies
	if err := client.DB.Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func GetMovieByTitle(ctx context.Context, title string) (Movies, error) {
	client, err := getNewDBClient(ctx)
	if err != nil {
		return Movies{}, err
	}

	var movie Movies
	if err := client.DB.Where("title = ?", title).First(&movie).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Movies{}, commonErrors.ErrRecordNotFound
		}
		return Movies{}, err
	}

	return movie, nil
}

func AddMovie(ctx context.Context, m Movies) error {
	client, err := getNewDBClient(ctx)
	if err != nil {
		return err
	}

	if err := client.DB.Create(&m).Error; err != nil {
		return err
	}

	return nil
}

func UpdateMovieLastFound(ctx context.Context, title string) error {
	client, err := getNewDBClient(ctx)
	if err != nil {
		return err
	}

	if err := client.DB.Model(&Movies{}).Where("title = ?", title).Updates(map[string]interface{}{
		"last_found_on": time.Now(),
		"occurrences":   gorm.Expr("occurrences + ?", 1),
	}).Error; err != nil {
		return err
	}

	return nil
}
