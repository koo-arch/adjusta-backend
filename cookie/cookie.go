package cookie

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/configs"
)

func init() {
	configs.LoadEnv()
}

// SetCookieはレスポンスにクッキーを設定します
func SetCookie(c *gin.Context, name, value string, maxAge int) {
	domain := configs.GetEnv("DOMAIN")
	println(domain)
	cookie:= &http.Cookie{
		Name: name,
		Value: value,
		MaxAge: maxAge,
		Path: "/",
		Domain: domain,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, cookie)
}

// DeleteCookie deletes a cookie in the response
func DeleteCookie(c *gin.Context, name string) {
	domain := configs.GetEnv("DOMAIN")
	cookie := &http.Cookie{
		Name: name,
		Value: "",
		MaxAge: -1,
		Path: "/",
		Domain: domain,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, cookie)
}