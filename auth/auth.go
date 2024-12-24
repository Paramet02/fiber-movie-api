package auth

import (
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	TokenExpiry  time.Duration
	RefreshExpiry time.Duration
	CookieName    string
	CookiePath    string
}

type jwtUser struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

type TokenPairs struct {
    Token        string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type Claims struct {
    jwt.RegisteredClaims
}

// ฟังก์ชันสร้าง jwtUser ใหม่
func NewJwtUser(id int, firstName, lastName string) *jwtUser {
    return &jwtUser{
        ID:        id,
        FirstName: firstName,
        LastName:  lastName,
    }
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (*TokenPairs, error) {

    // Create a token (สร้างโทเคน)
    token := jwt.New(jwt.SigningMethodHS256)

    // Set the claims (กำหนดข้อมูลเข้ารหัส)
    claims := token.Claims.(jwt.MapClaims)
    claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName) // ชื่อและนามสกุล
    claims["sub"] = fmt.Sprint(user.ID)                                  // รหัสผู้ใช้
    claims["aud"] = j.JWTAudience                                           // ผู้รับ JWT
    claims["iss"] = j.JWTIssuer                                             // ผู้ออก JWT
    claims["iat"] = time.Now().UTC().Unix()                              // วันที่และเวลาที่ออก JWT
    claims["typ"] = "JWT"                                                // ประเภทของ JWT

    // Set the expiry for JWT (กำหนดระยะเวลาในการใช้งาน JWT)
    claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

    // Create a signed token (สร้างโทเคนที่เข้ารหัสแล้ว)
    signedAccessToken, err := token.SignedString([]byte(j.JWTSecret))
    if err != nil {
        return &TokenPairs{}, err
    }

    // Create a refresh token and set claims (สร้าง Refresh Token และกำหนดข้อมูลเข้ารหัส)
    refreshToken := jwt.New(jwt.SigningMethodHS256)
    refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
    refreshTokenClaims["sub"] = fmt.Sprint(user.ID)     // รหัสผู้ใช้
    refreshTokenClaims["iat"] = time.Now().UTC().Unix() // วันที่และเวลาที่ออก Refresh Token

    // Set the expiry for the refresh token (กำหนดระยะเวลาในการใช้งาน Refresh Token)
    refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

    // Create signed refresh token (สร้าง Refresh Token ที่เข้ารหัสแล้ว)
    signedRefreshToken, err := refreshToken.SignedString([]byte(j.JWTSecret))
    if err != nil {
        return &TokenPairs{}, err
    }

    // Create TokenPairs and populate with signed tokens (สร้าง TokenPairs และเติมด้วยโทเคนที่เข้ารหัสแล้ว)
    // อยากให้ user มีสอง token เพราะเราจะดูว่าเคยเข้า web เรามั้ย ถ้าไม่ : ก็สร้างตัวใหม่ ถ้าใช่ : ก็ทำการ refresh token
    var tokenPairs = TokenPairs{
        Token:        signedAccessToken,
        RefreshToken: signedRefreshToken,
    }

    // Return TokenPairs and nil error (ส่งค่า TokenPairs และ nil ให้กับ error)
    return &tokenPairs, nil
}

// ฟังก์ชันสำหรับการ GetRefreshCookie
func(j *Auth) GetRefreshCookie(refreshToken string) *fiber.Cookie {
    return &fiber.Cookie{
        Name:     j.CookieName,
        Path:     j.CookiePath,
        Value:    refreshToken,
        Expires:  time.Now().Add(j.RefreshExpiry),
        MaxAge:   int(j.RefreshExpiry.Seconds()),
        SameSite: fiber.CookieSameSiteStrictMode,
        Domain:   j.CookieDomain,
        HTTPOnly: true,
        Secure:   true,
    }
}

// ฟังก์ชันสำหรับการ GetExpiredRefreshCookie
func (j *Auth) GetExpiredRefreshCookie() *fiber.Cookie {
    return &fiber.Cookie{
        Name:     j.CookieName,
        Path:     j.CookiePath,
        Value:    "",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
        SameSite: fiber.CookieSameSiteStrictMode,
        Domain:   j.CookieDomain,
        HTTPOnly: true,
        Secure:   true,
    }
}

// ฟังก์ชันสำหรับการ GetTokenFromHeaderAndVerify (Authorization Header)
func (j *Auth) GetTokenFromHeaderAndVerify(c *fiber.Ctx) (string, *Claims, error) {

    // get auth header
    authHeader := c.Get("Authorization")

    // sanity check
    if authHeader == "" {
        return "", nil, errors.New("no auth header")
    }

    // split the header on spaces
    headerParts := strings.Split(authHeader, " ")
    if len(headerParts) != 2 {
        return "", nil, errors.New("invalid auth header")
    }

    // check to see if we have the word Bearer
    if headerParts[0] != "Bearer" {
        return "", nil, errors.New("invalid auth header")
    }

    token := headerParts[1]

    // declare an empty claims
    claims := &Claims{}

    // parse the token
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(j.JWTSecret), nil
    })

    if err != nil {
        if strings.HasPrefix(err.Error(), "token is expired by") {
            return "", nil, errors.New("expired token")
        }
        return "", nil, err
    }

    if claims.Issuer != j.JWTIssuer {
        return "", nil, errors.New("invalid issuer")
    }

    return token, claims, nil
}

