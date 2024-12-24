package services

import (
	"errors"

	"github.com/paramet02/webapi/auth"
	"github.com/paramet02/webapi/repository"
    "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
    
)

type userService struct {
    userRepo repository.UserRepository
    auth     *auth.Auth
}

func NewuserService(userRepo repository.UserRepository, auth *auth.Auth) UserService {
    return &userService{userRepo, auth}
}

// Login function: Verifies user credentials and returns JWT tokens
func (s *userService) Login(email, password string) (*auth.TokenPairs, error) {
    user, err := s.userRepo.GetUserByEmail(email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Validate the password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil , errors.New("invalid credentials")
    }
    
    jwtUser := auth.NewJwtUser(user.ID , user.FirstName , user.LastName)


    return s.auth.GenerateTokenPair(jwtUser)
}

func (s *userService) Register(email, password, firstName, lastName string) (*auth.TokenPairs, error) {
    // Check if the email already exists
    existingUser, err := s.userRepo.GetUserByEmail(email)
    if err != nil && err != gorm.ErrRecordNotFound {
        return nil, err // return any error other than 'record not found'
    }
    
    if existingUser != nil {
        return nil, errors.New("email already in use")
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err // return error if password hashing fails
    }

    // Create a new User object with hashed password
    user := repository.NewUser(firstName, lastName, email, string(hashedPassword))

    // Save user to the database
    createdUser, err := s.userRepo.Create(user)
    if err != nil {
        return nil, err
    }

    // Create a JwtUser object for token generation
    jwtUser := auth.NewJwtUser(createdUser.ID, createdUser.FirstName, createdUser.LastName)

    // Generate and return the JWT token pair
    return s.auth.GenerateTokenPair(jwtUser)
}

func (s *userService) GetUserByEmail(Email string) (*User, error) {
    GetEmail, err := s.userRepo.GetUserByEmail(Email)
    if err != nil {
        return nil, err
    }

    user := &User{
        ID:        GetEmail.ID,
        FirstName: GetEmail.FirstName,
        LastName:  GetEmail.LastName,
        Email:     GetEmail.Email,
        Password:  GetEmail.Password,
        CreatedAt: GetEmail.CreatedAt,
        UpdatedAt: GetEmail.UpdatedAt,
    }

    return user, nil
}

func (s *userService) GetUserByID(ID int) (*User, error) {
    GetID, err := s.userRepo.GetUserByID(ID)
    if err != nil {
        return nil, err
    }

    user := &User{
        ID:        GetID.ID,
        FirstName: GetID.FirstName,
        LastName:  GetID.LastName,
        Email:     GetID.Email,
        Password:  GetID.Password,
        CreatedAt: GetID.CreatedAt,
        UpdatedAt: GetID.UpdatedAt,
    }

    return user, nil
}

