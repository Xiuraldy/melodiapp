package authcore

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	userports "melodiapp/internal/ports/user"
	"melodiapp/models"
	"melodiapp/shared"
)

type Service struct {
	repo userports.UserRepository
}

func NewService(repo userports.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(input models.UserInput) (string, error) {
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return "", errors.New("Incomplete fields")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(input.Email) {
		return "", errors.New("Invalid email format")
	}

	// Check existing by email
	existingByEmail, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return "", err
	}
	if existingByEmail != nil {
		return "", fmt.Errorf("Email already exists")
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := s.repo.CreateUser(&user); err != nil {
		// Duplicated username is handled via DB unique constraint
		return "", err
	}

	token, err := createSessionAndToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Logout(tokenStr string) error {
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	claims, _ := token.Claims.(*shared.Payload)
	_, exists := shared.Sessions[claims.Session]
	if !exists {
		return errors.New("You don't have permission")
	}

	delete(shared.Sessions, claims.Session)

	return nil
}

func (s *Service) Login(input models.UserInput) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", errors.New("Incomplete fields")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(input.Email) {
		return "", errors.New("Invalid email format")
	}

	user, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("Invalid credentials")
	}

	token, err := createSessionAndToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createSessionAndToken(userID uint) (string, error) {
	sessionToken := uuid.NewV5(uuid.UUID{}, "session").String()

	session := shared.Session{
		Uid:        userID,
		ExpiryTime: time.Now().Add(10 * time.Minute),
	}

	shared.Sessions[sessionToken] = session

	claims := shared.Payload{
		MapClaims: jwt.MapClaims{
			"iat":     jwt.NewNumericDate(time.Now()),
			"eat":     jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			"user_id": userID,
		},
		Session: sessionToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signinKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString, err := token.SignedString(signinKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
