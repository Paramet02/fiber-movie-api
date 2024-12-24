package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/paramet02/webapi/services"
)

// movieHandler โครงสร้างที่เชื่อมกับ Service
type movieHandler struct {
	MovieService services.MovieService
}

// NewMovieHandler สร้าง Handler และรับ Service เป็น Dependency
func NewmovieHandler(movieService services.MovieService) MovieHandler {
	return &movieHandler{MovieService: movieService}
}

func (h movieHandler) GetsMovies(c *fiber.Ctx) error {
    Movie, err := h.MovieService.GetsMovies()
    if err != nil {
        return c.SendStatus(fiber.StatusBadRequest)
    }

    // Convert the Handler (a slice of movies) into a JSON response
    return c.JSON(Movie)
}

func (h movieHandler) GetMovie(c *fiber.Ctx) error {
	id , err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	mHandler , err := h.MovieService.GetMovie(id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	
	return c.JSON(mHandler)
}

func (h movieHandler) OneMovieForEdit(c *fiber.Ctx) error {
	id , err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	movie , genre , err := h.MovieService.OneMovieForEdit(id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{
		"Movie" : movie,
		"Genre" : genre,
	})
}


func (h movieHandler) GetsGenres(c *fiber.Ctx) error {
	Genres, err := h.MovieService.GetsGenres()
    if err != nil {
        return c.SendStatus(fiber.StatusBadRequest)
    }

    // Convert the Handler (a slice of movies) into a JSON response
    return c.JSON(Genres)
}


func (h movieHandler) InsertMovie(c *fiber.Ctx) error {
    var movie services.Movie
    // ใช้ pointer ในการรับข้อมูลจาก request
    if err := c.BodyParser(&movie); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid input data",
        })
    }

    // เรียก getPoster เพื่อดึงข้อมูล poster
    movie, err := getPoster(movie)  // ส่ง pointer ไปที่ getPoster
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to fetch movie poster",
        })
    }

    // กำหนดเวลาให้กับ createdAt และ updatedAt
    movie.CreatedAt = time.Now()
    movie.UpdatedAt = time.Now()

    // ใช้ pointer ในการส่งข้อมูล
    MovieID, err := h.MovieService.InsertMovie(&movie)  // ส่ง pointer ไปที่ InsertMovie
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to insert movie 1",
        })
    }

    // อัพเดต Genres
    err = h.MovieService.UpdateMovieGenres(MovieID, movie.GenresArray)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to update movie genres",
        })
    }

    // ส่งข้อมูลกลับ
    return c.JSON(fiber.Map{
        "status":  "success",
        "movie_id": MovieID,
    })
}


func (h movieHandler) UpdateMovie(c *fiber.Ctx) error {
	var newMovie services.Movie
	id , err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	movie , err := h.MovieService.GetMovie(id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	movie.Title = newMovie.Title
	movie.ReleaseDate = newMovie.ReleaseDate
	movie.Description = newMovie.Description
	movie.MpaaRating = newMovie.MpaaRating
	movie.Runtime = newMovie.Runtime
	movie.UpdatedAt = time.Now()

	if err := h.MovieService.UpdateMovie(movie) ; err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.MovieService.UpdateMovieGenres(movie.ID , newMovie.GenresArray); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "success",
	})
}

func (h movieHandler) DeleteMovie(c *fiber.Ctx) error {
	id , err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = h.MovieService.DeleteMovie(id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "success",
	})
}

// ฟังก์ชันสำหรับดึงรูปภาพหนังจาก API
func getPoster(movie services.Movie) (services.Movie, error) {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}
	
	client := &http.Client{}
	APIKey := os.Getenv("API_KEY")
	theUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", APIKey)

	// สร้างคำขอสำหรับค้นหาหนังตามชื่อ
	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		log.Println(err)
		return movie, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// ส่งคำขอไปที่ API
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie, err
	}
	defer resp.Body.Close()

	// อ่านข้อมูลที่ได้รับจาก API
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie, err
	}

	var responseObject TheMovieDB
	// แปลงข้อมูล JSON เป็นโครงสร้างที่กำหนด
	if err := json.Unmarshal(bodyBytes, &responseObject); err != nil {
		log.Println(err)
		return movie, err
	}

	// ตรวจสอบว่ามีผลลัพธ์จากการค้นหาหนังหรือไม่
	if len(responseObject.Results) > 0 {
		movie.Image = responseObject.Results[0].PosterPath
	}

	// ส่งกลับ movie ที่ถูกอัปเดต
	return movie, nil
}

