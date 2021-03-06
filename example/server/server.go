package main

import (
	"log"
	"net/http"
	"github.com/hexagon-cloud/oauth2"
	"github.com/hexagon-cloud/oauth2/manager"
	"github.com/hexagon-cloud/oauth2/server"
	buntdbStore "github.com/hexagon-cloud/oauth2/store/buntdb"
	memoryStore "github.com/hexagon-cloud/oauth2/store/memory"
	"time"
	"github.com/hexagon-cloud/oauth2/password/hmac"
)

func main() {
	mgr := manager.NewDefaultManager()
	// token store
	mgr.MustTokenStorage(buntdbStore.NewMemoryTokenStore())

	// client store
	clientStore := memoryStore.NewClientStore()
	clientStore.Set("server", &oauth2.DefaultClient{
		ID:              "server",
		Secret:          "server",
		RedirectUri:     "http://localhost:8080",
		Scopes:          []string{"server", "all"},
		GrantTypes:      []oauth2.GrantType{oauth2.ClientCredentials},
		AccessTokenExp:  time.Duration(8) * time.Hour,
		RefreshTokenExp: time.Duration(8) * time.Hour,
	})
	clientStore.Set("app", &oauth2.DefaultClient{
		ID:              "app",
		Secret:          "app",
		Scopes:          []string{"app"},
		GrantTypes:      []oauth2.GrantType{oauth2.PasswordCredentials, oauth2.Implicit},
		AccessTokenExp:  time.Duration(8) * time.Hour,
		RefreshTokenExp: time.Duration(8) * time.Hour,
	})
	mgr.MapClientStorage(clientStore)

	// password encoder
	pwdEncoder := hmac.NewPasswordEncoder("key")
	mgr.MapPasswordEncoder(pwdEncoder)

	// user store
	userStore := memoryStore.NewUserStore()
	userStore.Set("user1", &oauth2.DefaultUser{
		ID:       1,
		Username: "user1",
		Password: pwdEncoder.Encode("pwd1"),
	})
	mgr.MapUserStorage(userStore)

	uaaServer := server.NewServer(server.NewConfig(), mgr)

	uaaServer.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		userID = "user1"
		return
	})
	uaaServer.SetInternalErrorHandler(func(err error) (re *oauth2.ErrorResponse) {
		log.Println("Internal Error:", err.Error())
		return
	})
	uaaServer.SetResponseErrorHandler(func(re *oauth2.ErrorResponse) {
		log.Println("ErrorResponse Error:", re.Error.Error())
	})

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := uaaServer.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		err := uaaServer.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/oauth/check_token", func(w http.ResponseWriter, r *http.Request) {
		err := uaaServer.HandleCheckTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/oauth/me", func(w http.ResponseWriter, r *http.Request) {
		err := uaaServer.HandleTokenUserRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server is running at 8401 port.")
	log.Fatal(http.ListenAndServe(":8401", nil))
}
