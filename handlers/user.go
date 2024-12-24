package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface{
	GetUserByEmail(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
} 