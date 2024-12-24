package repository

import (
	"time"
)

// movie โครงสร้างของข้อมูลหนัง
type movie struct {
	ID          int       `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(255)"`
	ReleaseDate time.Time `gorm:"type:date"`
	Runtime     int       `gorm:"not null"`
	MpaaRating  string    `gorm:"type:varchar(10)"`
	Description string    `gorm:"type:text"`
	Image       string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Genres []*genres `gorm:"many2many:movies_genres;" json:"genres,omitempty"`
	GenresArray []int `gorm:"-" json:"genres_array,omitempty"`
}

// genre โครงสร้างของข้อมูลประเภทหนัง
type genres struct {
    ID        int       `gorm:"primaryKey;autoIncrement"`
    Genre     string    `gorm:"type:varchar(255);not null"`
    Checked   bool      `gorm:"default:false"`
    CreatedAt time.Time
    UpdatedAt time.Time
	
    Movies    []*movie `gorm:"many2many:movies_genres;"`
}

type moviesgenres struct {
	ID 			int
	MovieId		int
	GenresId	int
}

// MovieRepository interface สำหรับ repository
type MovieRepository interface {
	AllMovies() ([]movie, error)         // Get all movies
	AllGenres() ([]genres, error)         // Get all movies
	GetMovie(id int) (*movie, error)     // Get a movie by ID
	OneMovieForEdit(id int) (*movie, []genres, error) // Update a movie by ID and return its genres
	InsertMovie(*movie) (int , error)
	UpdateMovie(*movie) error
	UpdateMovieGenres(movieID int, genreIDs []int) error
	DeleteMovie(id int) error
}

// Newmovie is a function that creates a new movie instance
func Newmovie(title string, releaseDate time.Time, runtime int, mpaaRating, description, image string) *movie {
	// สร้าง struct movie แล้วส่ง pointer กลับ
	return &movie{
		Title:       title,
		ReleaseDate: releaseDate,
		Runtime:     runtime,
		MpaaRating:  mpaaRating,
		Description: description,
		Image:       image,
	}
}

// NewUpdateMovie is a function used to update a movie's details by accepting the movie's id and the new values for name and quantity
func NewUpdateMovie(id int, title string, runtime int, mpaaRating string, description string, image string) *movie {
	// Create and return a movie struct with updated values based on the provided arguments
	return &movie{
		ID:          id,
		Title:       title,
		Runtime:     runtime,
		MpaaRating:  mpaaRating,
		Description: description,
		Image:       image,
	}
}
