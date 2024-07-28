package mw

import (
	"fast/config"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Conf.Https.Flag {
			prefix := strings.Split(c.Request.Host, ":")
			var host = "localhost"
			if len(prefix) >= 2 {
				host = prefix[0]
			}
			secureMiddleware := secure.New(secure.Options{
				SSLRedirect: true,
				SSLHost:     fmt.Sprintf("%s:%s", host, config.Conf.Https.Port),
			})
			err := secureMiddleware.Process(c.Writer, c.Request)
			if err != nil {
				return
			}
		}
		c.Next()
	}
}
