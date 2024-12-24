package main

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// authRequired ตรวจสอบการ Auth
func (app *Application) authRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, _ ,err := app.auth.GetTokenFromHeaderAndVerify(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}
		return c.Next()
	}
}

// Middleware สำหรับตรวจสอบ JWT Token
func (app *Application) jwtMiddleware(c *fiber.Ctx) error {
	// ดึง token จาก Authorization header
	tokenString := c.Get("Authorization")

	// ตรวจสอบว่ามีการส่ง Token มาหรือไม่
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// ตรวจสอบรูปแบบของ token ว่าเป็น Bearer token หรือไม่
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	// ตัดคำว่า "Bearer " ออก เพื่อให้เหลือเฉพาะ token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// ตรวจสอบความถูกต้องของ token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่ามีการใช้ signing method ที่ถูกต้องหรือไม่
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		// คืนค่า secret key ที่ใช้ในการตรวจสอบ
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// หาก token ถูกต้อง ให้ส่งคำขอไปยัง handler ถัดไป
	return c.Next()
}