package services

import (
	"log"

	"github.com/paramet02/webapi/repository"
)

// movieService คือตัว adapter สำหรับเชื่อมต่อกับ repository
type movieService struct {
	mService repository.MovieRepository
}

// NewMovieService สร้าง movieService โดยรับ MovieRepository เป็นพารามิเตอร์
func NewMovieService(mService repository.MovieRepository) MovieService {
	return &movieService{mService}
}

// GetsMovies ดึงข้อมูลหนังทั้งหมด
func (s *movieService) GetsMovies() ([]Movie, error) {
	movies, err := s.mService.AllMovies()
	if err != nil {
		return nil, err
	}

	var result []Movie
	for _, m := range movies {
		result = append(result, Movie{
			ID:          m.ID,
			Title:       m.Title,
			ReleaseDate: m.ReleaseDate,
			Runtime:     m.Runtime,
			MpaaRating:  m.MpaaRating,
			Description: m.Description,
			Image:       m.Image,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		})
	}
	return result, nil
}

// GetsMovies ดึงข้อมูลหนังทั้งหมด
func (s *movieService) GetsGenres() ([]Genres, error)  {
	movies, err := s.mService.AllGenres()
	if err != nil {
		return nil, err
	}

	var result []Genres
	for _, m := range movies {
		result = append(result, Genres{
			ID:    		m.ID,
			Genre: 		m.Genre,
			Checked:	m.Checked,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}
	return result, nil
}

// GetMovie ดึงข้อมูลหนังตาม ID
func (s *movieService) GetMovie(id int) (*Movie, error) {
	movie, err := s.mService.GetMovie(id)
	if err != nil {
		return nil, err
	}

	return &Movie{
		ID:          movie.ID,
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Runtime:     movie.Runtime,
		MpaaRating:  movie.MpaaRating,
		Description: movie.Description,
		Image:       movie.Image,
	}, nil
}

// UpdateMovie อัปเดตข้อมูลหนังและดึงข้อมูลประเภทหนัง
func (s *movieService) OneMovieForEdit(id int) (*Movie, []Genres, error) {
	movie, genres, err := s.mService.OneMovieForEdit(id)
	if err != nil {
		return nil, nil, err
	}

	var genreList []Genres
	for _, g := range genres {
		genreList = append(genreList, Genres{
			ID:      g.ID,
			Genre:   g.Genre,
			Checked: g.Checked,
		})
	}

	return &Movie{
		ID:          movie.ID,
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Runtime:     movie.Runtime,
		MpaaRating:  movie.MpaaRating,
		Description: movie.Description,
		Image:       movie.Image,
	}, genreList, nil
}

func (s *movieService) InsertMovie(movie *Movie) (int, error) {
    // สร้าง movie ใหม่เป็น pointer
    service := repository.Newmovie(movie.Title, movie.ReleaseDate, movie.Runtime, movie.MpaaRating, movie.Description, movie.Image)

    // ส่ง pointer ไปยัง InsertMovie
    movieID, err := s.mService.InsertMovie(service)
    if err != nil {
        log.Println("Failed to insert movie in service:", err)  // Log error
        return 0, err
    }

    // ถ้ามี GenresArray ก็จะต้อง update ความสัมพันธ์ Genres
    if len(movie.GenresArray) > 0 {
        err := s.mService.UpdateMovieGenres(movieID, movie.GenresArray)
        if err != nil {
            log.Println("Failed to update movie genres:", err)  // Log error
            return 0, err
        }
    }

    return movieID, nil
}


func (s *movieService) UpdateMovie(movie *Movie) error {
	// สร้าง movie ใหม่เป็น pointer
	service := repository.NewUpdateMovie(movie.ID , movie.Title , movie.Runtime , movie.MpaaRating , movie.Description , movie.Image)
	// ส่ง pointer ไปยัง UpdateMovie
	err := s.mService.UpdateMovie(service)
	if err != nil {
		return err
	}
	return nil
}


func (s *movieService) UpdateMovieGenres(movieID int, genreIDs []int) error {
	// เรียกใช้ UpdateMovieGenres จาก repository
	err := s.mService.UpdateMovieGenres(movieID, genreIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *movieService) DeleteMovie(id int) error {
	return s.mService.DeleteMovie(id)
}