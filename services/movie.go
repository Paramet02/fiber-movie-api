package services

import (
	"time"
)

// Movie โครงสร้างข้อมูลหนัง
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Runtime     int       `json:"runtime"`
	MpaaRating  string    `json:"mpaa_rating"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Genres []*Genres	  `json:"genres,omitempty"`
	GenresArray []int     `json:"genres_array,omitempty"`
}

// Genre โครงสร้างข้อมูลประเภทหนัง
type Genres struct {
	ID        int       `json:"id"`
	Genre     string    `json:"genre"`
	Checked   bool      `json:"checked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MovieService interface สำหรับ service
type MovieService interface {
	GetsMovies() ([]Movie, error)
	GetsGenres() ([]Genres, error) 
	GetMovie(id int) (*Movie, error)
	OneMovieForEdit(id int) (*Movie, []Genres, error)
	InsertMovie(*Movie) (int , error)
	UpdateMovie(*Movie) error
	UpdateMovieGenres(movieID int, genreIDs []int) error
	DeleteMovie(id int) error
}
