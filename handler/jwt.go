package handler

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/dgrijalva/jwt-go"
	"github.com/dikaeinstein/gomicroservice-auth/config"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// customCliams represents the cliams used to create a JWT token
type customCliams struct {
	UserID      string `json:"userId"`
	AccessLevel string `json:"accessLevel"`
	jwt.StandardClaims
}

// JWT is an http handler for the auth service
type JWT struct {
	logger     *log.Logger
	statsd     *statsd.Client
	rsaPrivate *rsa.PrivateKey
}

// loginRequest defines a request sent to the /login endpoint
type loginRequest struct {
	Username string `json:"username" validate:"email"`
	Password string `json:"password" validate:"max=36,min=8"`
}

// NewJWT creates a new JWT handler
func NewJWT(l *log.Logger, s *statsd.Client, cfg config.Config) *JWT {
	b, err := config.ParseRsaPrivateKeyHex(cfg.RsaPrivateKey)
	if err != nil {
		log.WithFields(defaultFields).
			Fatalf("Unable to decode private key hex: %v", err)
	}

	rsaPrivate, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		log.WithFields(defaultFields).
			Fatalf("Unable to parse private key: %v", err)
	}

	return &JWT{l, s, rsaPrivate}
}

func (j *JWT) generateJWT(request loginRequest) []byte {
	cliams := customCliams{
		"dikaeinstein",
		"user",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "dika",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, cliams)

	tokenString, err := token.SignedString(j.rsaPrivate)
	if err != nil {
		panic(err)
	}

	return []byte(tokenString)
}

var defaultFields = log.Fields{
	"service": "auth",
	"handler": "jwt",
}

// Login handler
func (j *JWT) Login(w http.ResponseWriter, r *http.Request) {
	defer func(startTime time.Time) {
		j.statsd.Timing("jwt.timing", time.Since(startTime), nil, 1)
	}(time.Now())

	request := loginRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		j.statsd.Incr("jwt.badrequest", nil, 1)

		j.logger.WithFields(defaultFields).
			Errorf("Error decoding request: %v", err)
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		j.statsd.Incr("jwt.badrequest", nil, 1)

		j.logger.WithFields(defaultFields).
			Errorf("Error validating request: %v", err)
		http.Error(w, "Error validating request", http.StatusBadRequest)
		return
	}

	j.logger.WithFields(defaultFields).
		Infof("Login request from %s", request.Username)
	jwt := j.generateJWT(request)

	w.Write(jwt)
	j.statsd.Incr("jwt.success", nil, 1)
}
