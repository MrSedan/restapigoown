package apiserver

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MrSedan/restapigoown/backend/internal/app/model"
	"github.com/MrSedan/restapigoown/backend/internal/app/store"
	"github.com/MrSedan/restapigoown/backend/internal/app/websockets"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	wsServ *websockets.Server
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
		wsServ: websockets.NewServer(store),
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	s.configureRouter()
	go s.wsServ.Run()
	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/", s.handleHome())
	s.router.HandleFunc(`/chat/{id:[0-9]+\.[0-9]+}`, s.handleWs())
	s.router.Handle(`/chat/{id:[0-9]+\.[0-9]+}/gethistory`, s.auth(s.handleGetMessagesHistory())).Methods("POST")
	s.router.Handle("/user/getalluser", s.auth(s.handleGetAllUser())).Methods("POST")
	s.router.Handle("/checkauth", s.auth(s.handleCheckAuth())).Methods("POST")
	s.router.HandleFunc("/user/create", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/user/login", s.handleLoginUser()).Methods("POST")
	s.router.HandleFunc("/user/{id:[0-9]+}/profile", s.handleProfile()).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/avatar", s.handleGetUserAvatar()).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/edit/profile", s.handleEditProfile()).Methods("POST")
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

func (s *server) handleGetAllUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.User().GetAllUsers()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		if err := json.NewEncoder(w).Encode(users); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
	}
}

func (s *server) handleGetMessagesHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PostFormValue("id")
		id := mux.Vars(r)["id"]
		ids := strings.Split(id, ".")
		if len(ids) != 2 {
			s.error(w, r, http.StatusBadRequest, errors.New("Not valid id"))
			return
		}
		p1, _ := strconv.Atoi(ids[0])
		p2, _ := strconv.Atoi(ids[1])
		if !(ids[0] == userID || ids[1] == userID) {
			s.error(w, r, http.StatusBadRequest, errors.New("Is not a your chat"))
			return
		}
		messages, err := s.store.User().GetMessageHistory(p1, p2)
		if err == store.ErrNotMessages {
			s.error(w, r, http.StatusNoContent, err)
			return
		} else if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Time < messages[j].Time
		})
		if err := json.NewEncoder(w).Encode(messages); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

	}
}

func (s *server) handleCheckAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
}

func (s *server) handleWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = newResponseWriter(w)
		id := mux.Vars(r)["id"]
		ids := strings.Split(id, ".")
		if len(ids) != 2 {
			s.error(w, r, http.StatusBadRequest, errors.New("Not valid id"))
			return
		}
		myid := ids[0]
		sort.Strings(ids)
		id = strings.Join(ids, ".")
		hub, ok := s.wsServ.Hubs[id]
		if !ok {
			hub = websockets.NewHub(id, s.wsServ)
			s.wsServ.NewHub <- hub
			go hub.Run()
		}
		websockets.ServeWs(hub, myid, w, r)
	}
}

func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		UserName string `json:"user_name"`
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
			UserName: req.UserName,
			Email:    strings.ToLower(req.Email),
			Password: req.Password,
		}
		if _, err := s.store.User().FindByEmail(u.Email); err == nil {
			s.error(w, r, http.StatusBadRequest, errors.New("Dublicate account email"))
			return
		}
		if _, err := s.store.User().FindByNick(u.UserName); err == nil {
			s.error(w, r, http.StatusBadRequest, errors.New("Dublicate account username"))
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
		token := jwt.New(jwt.SigningMethodHS256)
		rtClaims := token.Claims.(jwt.MapClaims)
		rtClaims["sub"] = 1
		rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(s.jwtKey)
		s.store.User().ClaimToken(u, tokenString)
		tokenString = fmt.Sprint(tokenString)
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

func (s *server) handleEditProfile() http.HandlerFunc {
	type request struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Token     string `json:"token"`
		About     string `json:"about,omitempty"`
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
		tokenID, err := s.store.User().CheckToken(bearerToken)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		u, err := s.store.User().FindByID(em)
		if err != nil || em != tokenID {
			s.error(w, r, http.StatusBadRequest, errUserErr)
			return
		}
		about := req.About
		firstName := req.FirstName
		lastName := req.LastName
		if about == "" && firstName == "" && lastName == "" {
			s.error(w, r, http.StatusBadRequest, errNotAboutField)
			return
		}
		if about != "" {
			if err := s.store.User().EditAbout(u.ID, about); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}
		if firstName != "" {
			if err := s.store.User().EditFirstName(u.ID, firstName); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}
		if lastName != "" {
			if err := s.store.User().EditLastName(u.ID, lastName); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
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
		tokenString, err := s.store.User().GetToken(u.ID)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, store.ErrNotValidToken)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"id": u.ID, "token": tokenString})
	}
}

func (s *server) handleProfile() http.HandlerFunc {
	type response struct {
		ID        int    `json:"id"`
		UserName  string `json:"user_name"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		About     string `json:"about,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := mux.Vars(r)["id"]
		u, err := s.store.User().FindByID(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		up, err := s.store.User().GetProfile(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		res := &response{
			ID:        u.ID,
			UserName:  u.UserName,
			FirstName: up.FirstName,
			LastName:  up.LastName,
			About:     up.About,
		}
		s.respond(w, r, http.StatusOK, res)
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
		tokenID, err := s.store.User().CheckToken(bearerToken)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		u, err := s.store.User().FindByID(em)
		if err != nil || string(u.ID) != tokenID || !u.ComparePassword(req.OldPass) {
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
		rw := newResponseWriter(w)
		rw.code = http.StatusOK
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		token := r.PostFormValue("token")
		id := r.PostFormValue("id")
		idUsr, err := s.store.User().CheckToken(token)
		if err != nil || token == "" || id != idUsr {
			wr := newResponseWriter(w)
			wr.code = http.StatusUnauthorized
			http.Error(wr, http.StatusText(wr.code), wr.code)
			return
		}
		rw := newResponseWriter(w)
		rw.code = http.StatusOK
		next.ServeHTTP(rw, r)
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
