package apiserver

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/MrSedan/restapigoown/backend/internal/app/model"
	"github.com/MrSedan/restapigoown/backend/internal/app/store"
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
	errNotAboutField            = errors.New("not about field")
	errUserErr                  = errors.New("incorrect email or password")
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
	s.router.HandleFunc("/user/create", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/user/login", s.handleLoginUser()).Methods("POST")
	s.router.HandleFunc("/user/{id:[0-9]+}/profile", s.handleProfile()).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/avatar", s.handleGetUserAvatar()).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/edit/profile", s.handleEditAbout()).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/edit/password", s.handleEditPassword()).Methods("POST")
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
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     strings.ToLower(req.Email),
			Password:  req.Password,
		}
		if _, err := s.store.User().FindByEmail(u.Email); err == nil {
			s.error(w, r, http.StatusBadRequest, errors.New("dublicate account"))
			return
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

func (s *server) handleGetUserAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		id := mux.Vars(r)["id"]
		u, err := s.store.User().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		resp, err := http.Get(fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=300&d=identicon", md5.Sum([]byte(u.Email))))
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	}
}

func (s *server) handleEditAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		bearerToken := r.FormValue("token")
		em := mux.Vars(r)["id"]
		tokenEmail, _ := s.store.User().GetToken(bearerToken)
		u, err := s.store.User().FindByID(em)
		if err != nil || em != tokenEmail {
			s.error(w, r, http.StatusBadRequest, errUserErr)
			return
		}
		about := r.FormValue("about")
		if about == "" {
			s.error(w, r, http.StatusBadRequest, errNotAboutField)
			return
		}
		if err := s.store.User().EditAbout(u.ID, about); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]string{"status": "ok"})
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
		u, err := s.store.User().FindByEmail(strings.ToLower(req.Email))
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
		tokenString = fmt.Sprint(tokenString)
		s.respond(w, r, http.StatusOK, map[string]interface{}{"id": u.ID, "token": tokenString})
	}
}

func (s *server) handleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := mux.Vars(r)["id"]
		u, err := s.store.User().GetProfile(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleEditPassword() http.HandlerFunc {
	type request struct {
		OldPass string `json:"old_password"`
		NewPass string `json:"new_password"`
		Token   string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		bearerToken := req.Token
		em := mux.Vars(r)["id"]
		tokenID, _ := s.store.User().GetToken(bearerToken)
		u, err := s.store.User().FindByID(em)
		if err != nil || em != tokenID || !u.ComparePassword(req.OldPass) {
			s.error(w, r, http.StatusBadRequest, errUserErr)
			return
		}
		u.Password = req.NewPass
		if err := s.store.User().EditPass(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]string{"status": "ok"})
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
