package handlers

import (
	"strconv"
    "errors"

	"github.com/gofiber/fiber/v2"
	"github.com/paramet02/webapi/services"
    "github.com/paramet02/webapi/auth"
    "github.com/golang-jwt/jwt/v4"
)

// Adapter: โครงสร้างที่ทำหน้าที่ implement ฟังก์ชันใน Interface (Port) โดยเชื่อมต่อกับฐานข้อมูล
type userHandler struct {
	userService services.UserService
    auth     *auth.Auth
}

// ฟังก์ชันสำหรับสร้าง MovieRepository ด้วยการรวม Port (Interface) และ Adapter (Implementation)
func NewuserHandler(userHand services.UserService ,auth *auth.Auth) UserHandler {
	return &userHandler{userHand , auth}
}

func (h userHandler) GetUserByEmail(c *fiber.Ctx) error {
	// Query Parameters
	Email := c.Query("Email") 
	
	// ถ้าไม่มี email ให้ส่งสถานะ 400
    if Email == "" {
        return c.SendStatus(fiber.StatusBadRequest)
    }

	// รับค่าจาก service
	userEmail , err := h.userService.GetUserByEmail(Email)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// ส่งผลลัพธ์กลับไปในรูปแบบ JSON
    return c.JSON(fiber.Map{
        "status":  "ok",
        "products": userEmail, // ส่งสินค้าในรูปแบบ JSON
    })
}

func (h userHandler) GetUserByID(c *fiber.Ctx) error {
	id , err := strconv.Atoi(c.Params("id"))

	// ถ้า id ผิด
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID , err := h.userService.GetUserByID(id)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
    
	// ส่งผลลัพธ์กลับไปในรูปแบบ JSON
    return c.JSON(fiber.Map{
        "status":  "ok",
        "products": userID, // ส่งสินค้าในรูปแบบ JSON
    })
}
// Login endpoint
func (h userHandler) Login(c *fiber.Ctx) error {
    var request struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Parse the JSON body
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "error",
            "message": "Invalid request body",
        })
    }

    tokens, err := h.userService.Login(request.Email, request.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "status": "error",
            "message": err.Error(),
        })
    }

     // Get the user that was just created using email
     createdUser, err := h.userService.GetUserByEmail(request.Email)
     if err != nil {
         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
             "status":  "error",
             "message": "Failed to retrieve user after registration",
         })
     }
 
     // Set refresh token cookie
     refreshCookie := h.auth.GetRefreshCookie(tokens.RefreshToken)
     c.Cookie(&fiber.Cookie{
         Name:     refreshCookie.Name,
         Value:    refreshCookie.Value,
         Expires:  refreshCookie.Expires,
         Path:     refreshCookie.Path,
         Secure:   refreshCookie.Secure,
         HTTPOnly: refreshCookie.HTTPOnly,
     })
 
     // Create the response payload
     responsePayload := struct {
         AccessToken  string `json:"access_token"`
         RefreshToken string `json:"refresh_token"`
         User         struct {
             ID        int    `json:"id"`
             FirstName string `json:"first_name"`
             LastName  string `json:"last_name"`
             Email     string `json:"email"`
         } `json:"user"`
     }{
         AccessToken:  tokens.Token,
         RefreshToken: tokens.RefreshToken,
         User: struct {
             ID        int    `json:"id"`
             FirstName string `json:"first_name"`
             LastName  string `json:"last_name"`
             Email     string `json:"email"`
         }{
             ID:        createdUser.ID,  // Use createdUser's ID
             FirstName: createdUser.FirstName,
             LastName:  createdUser.LastName,
             Email:     createdUser.Email,
         },
     }
 
     // Write the response as JSON
     return c.Status(fiber.StatusAccepted).JSON(responsePayload)
 }

// Register endpoint
func (h userHandler) Register(c *fiber.Ctx) error {
    var request struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
    }

    // Parse the JSON body
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "error",
            "message": "Invalid request body",
        })
    }

    tokens, err := h.userService.Register(request.Email, request.Password, request.FirstName, request.LastName)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "error",
            "message": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "status": "ok",
        "tokens": tokens,
    })
}

func (h userHandler) Logout(c *fiber.Ctx) error {
	// เรียกใช้ GetExpiredRefreshCookie() และตั้งคุกกี้ที่หมดอายุ
	c.Cookie(h.auth.GetExpiredRefreshCookie())
	
	// ส่งสถานะ HTTP 202 Accepted
	return c.SendStatus(fiber.StatusAccepted)
}


func (h userHandler) RefreshToken(c *fiber.Ctx) error {
	// Get the refresh token from the cookies
	refreshTokenCookie := c.Cookies(h.auth.CookieName)

	// If the cookie is not present, return an unauthorized error
	if refreshTokenCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Refresh token missing",
		})
	}

	// Parse the refresh token to get the claims
	claims := &auth.Claims{}
	_, err := jwt.ParseWithClaims(refreshTokenCookie, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.auth.JWTSecret), nil
	})

	// If there is an error parsing the token, return unauthorized
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid refresh token",
		})
	}

	// Get the user ID from the claims
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID in token",
		})
	}

	// Retrieve the user from the database
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	// Create a new JwtUser
	jwtUser := auth.NewJwtUser(user.ID, user.FirstName, user.LastName)

	// Generate new tokens
	tokens, err := h.auth.GenerateTokenPair(jwtUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Error generating tokens",
		})
	}

	// Set the new refresh token cookie
	refreshCookie := h.auth.GetRefreshCookie(tokens.RefreshToken)
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookie.Name,
		Value:    refreshCookie.Value,
		Expires:  refreshCookie.Expires,
		Path:     refreshCookie.Path,
		Secure:   refreshCookie.Secure,
		HTTPOnly: refreshCookie.HTTPOnly,
	})

	// Return the new tokens in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"tokens": tokens,
	})
}

