package middleware

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const secret = "csrf_secret"

func CSRF() gin.HandlerFunc {

	getToken := func(tsStr string) string {
		h := sha1.New()
		io.WriteString(h, tsStr+"-"+secret)
		token := base64.URLEncoding.EncodeToString(h.Sum(nil)) + "|" + tsStr
		return token
	}

	return func(c *gin.Context) {
		ts := time.Now().UnixNano() / int64(time.Millisecond)
		tsStr := strconv.FormatInt(ts, 10)
		token := getToken(tsStr)
		c.Set("csrfToken", token)
		if c.Request.Method == "POST" {
			token := c.Request.FormValue("csrfToken")
			if token == "" {
				panic("token not present")
			}
			tokenParts := strings.Split(token, "|")
			if len(tokenParts) != 2 {
				panic("bad token")
			}
			tsStr := tokenParts[1]
			calcToken := getToken(tsStr)
			if token != calcToken {
				panic("bad token!")
			}
		}
		c.Next()
	}
}
