package helper

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type CookieStore struct {
	Store map[int]string
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request, cookieStore *CookieStore) {
	cookie := getCookie(w, r)
	fmt.Println("Get cookie handler hit")
	fmt.Printf("Cookie: %v\n", cookie)
	if cookie != "" {
		for key, val := range cookieStore.Store {
			fmt.Printf("Key: %v, val: %v\n", key, val)
			if val == cookie {
				fmt.Printf("Welcome back user: %v\n", key)
			}
		}
	}
}

func getCookie(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("regCookie")
	if err != nil {
		switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusBadRequest)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
		}

		return ""
	}

	return cookie.Value
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request, uid string, cookieStore *CookieStore) {
	// setCookie(w, uuid.New(), cookieStore)
	cookie := http.Cookie{
		Name: "regCookie",
		Value: uid,
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	fmt.Printf("Set cookie %+v\n", cookie)

	cookieStore.Store[len(cookieStore.Store)+1] = cookie.Value

	http.SetCookie(w, &cookie)
	fmt.Println("Cookie handler set completed")
}