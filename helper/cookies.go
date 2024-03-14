package helper

import (
	"github.com/gorilla/securecookie"
	"net/http"
	//"fmt"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	cookieHandler *securecookie.SecureCookie
)

// Initialize the secure cookie encoder and decoder
func init() {
	godotenv.Load()
	key64, _ := hex.DecodeString(os.Getenv("KEY64"))
	key32, _ := hex.DecodeString(os.Getenv("KEY32"))
	cookieHandler = securecookie.New(
		key64,
		key32,
	)
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request, token string, refresh_token string) {
	value := map[string]string{
		"token":         token,
		"refresh_token": refresh_token,
	}

	if encoded, err := cookieHandler.Encode("cookie-name", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie-name",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
	if _, err := cookieHandler.Encode("cookie-name", value); err != nil {
		fmt.Println(err)
	}

}

func ReadCookie(r *http.Request) (map[string]string, error) {
	// Retrieve the cookie from the request
	cookie, err := r.Cookie("cookie-name")
	if err != nil {
		return nil, err
	}

	// Create a map to store the decoded cookie value
	value := make(map[string]string)

	// Decode the cookie value
	if err = cookieHandler.Decode("cookie-name", cookie.Value, &value); err != nil {
		return nil, err
	}

	return value, nil
}
