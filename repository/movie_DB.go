package repository

import (
	"log"

	"gorm.io/gorm"
)

// movieRepositoryDB คือตัว adapter สำหรับเชื่อมต่อกับฐานข้อมูล
type movieRepositoryDB struct {
	db *gorm.DB
}

// NewMovieRepositoryDB สร้าง MovieRepository โดยรับฐานข้อมูลเป็นพารามิเตอร์
func NewMovieRepositoryDB(db *gorm.DB) MovieRepository {
	return &movieRepositoryDB{db}
}

// AllMovies ดึงข้อมูลทั้งหมดของหนัง
func (m *movieRepositoryDB) AllMovies() ([]movie, error) {
	var movies []movie
	if result := m.db.Find(&movies); result.Error != nil {
		return nil, result.Error
	}
	return movies, nil
}

// GetMovie ดึงข้อมูลหนังตาม ID
func (m *movieRepositoryDB) GetMovie(id int) (*movie, error) {
	var movie movie
	if result := m.db.Preload("Genres").First(&movie, id); result.Error != nil {
		return nil, result.Error
	}
	return &movie, nil
}

// UpdateMovie อัปเดตข้อมูลหนังและดึงข้อมูลประเภทหนังทั้งหมด
func (m *movieRepositoryDB) OneMovieForEdit(id int) (*movie, []genres, error) {
	var resultMovie movie
	err := m.db.Preload("Genres").First(&resultMovie, id).Error
	if err != nil {
		return nil, nil, err
	}
	
	var allGenres []genres
	err = m.db.Order("genre").Find(&allGenres).Error
	if err != nil {
		return nil, nil, err
	}

	return &resultMovie, allGenres, nil
}

func (m *movieRepositoryDB) AllGenres() (genres []genres, err error) {
	if result := m.db.Find(&genres); result.Error != nil {
		return nil, result.Error
	}
	return genres , nil
}

func (m *movieRepositoryDB) InsertMovie(movie *movie) (int, error) {
    // ตรวจสอบค่าของ movie.ID หลังจาก insert
    result := m.db.Create(movie)
    if result.Error != nil {
        log.Println("Error inserting movie:", result.Error)  // Log error
        return 0, result.Error
    }

    return movie.ID, nil
}


func (m *movieRepositoryDB) UpdateMovie(movie *movie) error {
	if err := m.db.Model(&movie).Where("id = ? ", movie.ID).Updates(map[string]interface{}{
		"Title" : movie.Title,
		"ReleaseDate" : movie.ReleaseDate,
		"Runtime" : movie.Runtime,
		"MpaaRating" : movie.MpaaRating,
		"Description" : movie.Description,
		"Image" : movie.Image,
	}).Error; err != nil {
		return nil
	}
	
	return nil
}

func (m *movieRepositoryDB) UpdateMovieGenres(movieID int, genreIDs []int) error {
    // ลบ genres เดิมออก
    if err := m.db.Where("movieID = ?", movieID).Delete(&moviesgenres{}).Error ; err != nil {
		return nil
	}

	// เตรียมข้อมูลใหม่สำหรับการ insert
	var movieGenres []moviesgenres
	for _, genreID := range genreIDs {
		movieGenres = append(movieGenres, moviesgenres{
			MovieId: movieID,
			GenresId: genreID,
		})
	}

	// เพิ่ม genres ใหม่
	if err := m.db.Create(&movieGenres).Error; err != nil {
		return err
	}

    return nil
}

func (m *movieRepositoryDB) DeleteMovie(id int) error {
	if err := m.db.Where("id = ?" , id).Delete(&movie{}).Error; err != nil {
		return err
	}
	
	return nil 
}
