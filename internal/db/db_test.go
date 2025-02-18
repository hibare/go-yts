package db

import (
	"context"
	"os"
	"testing"

	commonErrors "github.com/hibare/GoCommon/v2/pkg/errors"
	"github.com/hibare/go-yts/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetDSN(t *testing.T) {
	os.Setenv("GO_YTS_DATA_DIR", "/test/data/dir")
	defer os.Clearenv()

	config.LoadConfig()

	dsn := getDSN()
	assert.Equal(t, "/test/data/dir/movies.db", dsn)
}

func TestMovie(t *testing.T) {
	ctx := context.Background()

	var oldGetDSN = getDSN
	getDSN = func() string {
		return ":memory:"
	}

	t.Cleanup(func() {
		getDSN = oldGetDSN
	})

	// Migrate db
	c, err := getNewDBClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Migrate(); err != nil {
		t.Fatal(err)
	}

	t.Run("AddMovie", func(t *testing.T) {
		actualMovie := Movies{
			Title:       "Test Movie",
			Link:        "https://example.com",
			Year:        2021,
			Occurrences: 1,
		}
		err := AddMovie(context.Background(), actualMovie)

		assert.NoError(t, err)

		var expectedMovie Movies
		err = c.DB.Where("title = ?", "Test Movie").First(&expectedMovie).Error
		assert.NoError(t, err)
		assert.Equal(t, actualMovie.Title, expectedMovie.Title)
		assert.Equal(t, actualMovie.Link, expectedMovie.Link)
		assert.Equal(t, actualMovie.Year, expectedMovie.Year)
		assert.Equal(t, actualMovie.Occurrences, expectedMovie.Occurrences)
		assert.NotEmpty(t, expectedMovie.FirstFoundOn)
	})

	t.Run("GetMovieByTitle", func(t *testing.T) {
		t.Run("Found", func(t *testing.T) {
			movie, err := GetMovieByTitle(ctx, "Test Movie")
			assert.NoError(t, err)
			assert.Equal(t, "Test Movie", movie.Title)
		})

		t.Run("NotFound", func(t *testing.T) {
			_, err := GetMovieByTitle(ctx, "Non Existent Movie")
			assert.Error(t, err)
			assert.ErrorIs(t, err, commonErrors.ErrRecordNotFound)
		})
	})

	t.Run("UpdateMovieLastFound", func(t *testing.T) {
		err := UpdateMovieLastFound(ctx, "Test Movie")
		assert.NoError(t, err)

		var movie Movies
		err = c.DB.Where("title = ?", "Test Movie").First(&movie).Error
		assert.NoError(t, err)
		assert.NotEqual(t, movie.LastFoundOn, movie.FirstFoundOn)
		assert.Greater(t, movie.LastFoundOn, movie.FirstFoundOn)
	})

	t.Run("GetAllMovies", func(t *testing.T) {
		movies, err := GetAllMovies(ctx)
		assert.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Title)
	})

}
