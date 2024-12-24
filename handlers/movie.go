package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type MovieHandler interface {
	GetsMovies(c *fiber.Ctx) error
	GetMovie(c *fiber.Ctx) error
	OneMovieForEdit(c *fiber.Ctx) error
	GetsGenres(c *fiber.Ctx) error
	InsertMovie(c *fiber.Ctx) error
	UpdateMovie(c *fiber.Ctx) error
	DeleteMovie(c *fiber.Ctx) error
}