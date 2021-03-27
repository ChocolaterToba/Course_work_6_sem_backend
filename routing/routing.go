package routing

import (
	"log"
	"net/http"
	"pinterest/auth"
	"pinterest/pins"
	"pinterest/profile"

	"github.com/gorilla/mux"
)

func boardHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{}"))
	// TODO /pin and /pins handling
}

// PanicMid logges error if handler errors
func PanicMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(PanicMid)

	r.HandleFunc("/auth/signup", auth.NoAuthMid(auth.HandleCreateUser)).Methods("POST")
	r.HandleFunc("/auth/login", auth.NoAuthMid(auth.HandleLoginUser)).Methods("POST")
	r.HandleFunc("/auth/logout", auth.AuthMid(auth.HandleLogoutUser)).Methods("POST")
	r.HandleFunc("/auth/check", auth.HandleCheckUser).Methods("GET")

	r.HandleFunc("/profile/password", auth.AuthMid(profile.HandleChangePassword)).Methods("PUT")
	r.HandleFunc("/profile/edit", auth.AuthMid(profile.HandleEditProfile)).Methods("PUT")
	r.HandleFunc("/profile/delete", auth.AuthMid(profile.HandleDeleteProfile)).Methods("DELETE")
	r.HandleFunc("/profile/{id:[0-9]+}", profile.HandleGetProfile).Methods("GET") // Is preferred over next one
	r.HandleFunc("/profile/{username}", profile.HandleGetProfile).Methods("GET")
	r.HandleFunc("/profile", auth.AuthMid(profile.HandleGetProfile)).Methods("GET")

	pins := &pins.PinsStorage{
		Storage: pins.NewPinsSet(),
	}

	r.HandleFunc("/pin", auth.AuthMid(pins.Storage.AddPin)).Methods("POST")
	r.HandleFunc("/pin/{id:[0-9]+}", pins.Storage.GetPinByID).Methods("GET")
	r.HandleFunc("/pin/{id:[0-9]+}", auth.AuthMid(pins.Storage.DelPinByID)).Methods("DELETE")
	r.HandleFunc("/pins/{id:[0-9]+}", auth.AuthMid(pins.Storage.DelPinByID)).Methods("DELETE")
	r.HandleFunc("/board/", boardHandler) // Will split later

	return r
}
