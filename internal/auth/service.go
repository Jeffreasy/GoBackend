package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Jeffreasy/GoBackend/configs"
	"github.com/Jeffreasy/GoBackend/internal/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(user *models.User) error
	Authenticate(email, password string) (string, error)
}

type authService struct {
	db  *sql.DB
	cfg *configs.Config
}

func NewService(db *sql.DB, cfg *configs.Config) Service {
	return &authService{
		db:  db,
		cfg: cfg,
	}
}

func (s *authService) RegisterUser(user *models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3)`
	_, err = s.db.Exec(query, user.Email, string(hashed), user.Name)
	return err
}

func (s *authService) Authenticate(email, password string) (string, error) {
	var hashedPwd string
	var id int
	query := `SELECT id, password FROM users WHERE email=$1`
	err := s.db.QueryRow(query, email).Scan(&id, &hashedPwd)
	if err != nil {
		return "", errors.New("ongeldige inloggegevens")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password)); err != nil {
		return "", errors.New("ongeldige inloggegevens")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
