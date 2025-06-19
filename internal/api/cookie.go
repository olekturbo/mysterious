package api

import "net/http"

const idCookie = "userId"

type CookieManager struct{}

func NewCookieManager() *CookieManager {
	return &CookieManager{}
}

func (c *CookieManager) Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (c *CookieManager) Write(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
}
