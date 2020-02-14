package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MrSedan/restapigoown/internal/app/model"
	"github.com/MrSedan/restapigoown/internal/app/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
	jwtKey []byte
}

const (
	ctxKeyRequestID = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/", s.handleHome())
	s.router.HandleFunc("/createuser", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLoginUser()).Methods("POST")
	s.router.Handle("/myprofile", s.validateToken(s.handleMyProfile())).Methods("POST")
	s.router.Handle("/test", s.validateToken(s.handleHome()))
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Handlers
func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
}

func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.store.User().CreateProfile(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleLoginUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)
		rtClaims := token.Claims.(jwt.MapClaims)
		rtClaims["sub"] = 1
		rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(s.jwtKey)
		s.store.User().ClaimToken(u, tokenString)
		tokenString = fmt.Sprintf("Bearer %s", tokenString)
		w.Header().Add("Authorization", tokenString)
		s.respond(w, r, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func (s *server) handleMyProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := r.Header.Get("authorization")
		bearerToken := strings.Split(req, " ")
		email, err := s.store.User().GetToken(bearerToken[1])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		prof := s.store.User().GetProfile(email)
		s.respond(w, r, http.StatusOK, prof)
	}
}

// MiddleWares
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Header.Get("authorization")
		if req != "" {
			bearerToken := strings.Split(req, " ")
			if len(bearerToken) == 2 {
				if _, err := s.store.User().GetToken(bearerToken[1]); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				rw := &responseWriter{w, http.StatusOK}
				next.ServeHTTP(rw, r)
			} else {
				s.error(w, r, http.StatusUnprocessableEntity, store.ErrNotValidToken)
				return
			}

		} else {
			s.error(w, r, http.StatusBadRequest, errors.New("no token"))
			return
		}
	})
}

// Helpers
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
