// token provides JWT validation and parsing.
package token

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/deuscapturus/tism/config"
	"github.com/deuscapturus/tism/randid"
	"github.com/deuscapturus/tism/request"
	"github.com/deuscapturus/tism/utils"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"time"
)

type JwtClaimsMap struct {
	Keys       []string `json:"keys"`
	Admin      int      `json:"admin"`
	JWTid      string   `json:"jti"`
	Expiration int64    `json:"exp"`
	jwt.StandardClaims
}

// ValidateJWT validate string jwt and return true/false if valid along with a slice of uint64 pgp key id's.

type RequestDecrypt struct {
	Token string `json:"token"`
}

func parseToken(t string) (token *jwt.Token, err error) {
	signingSecret := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWTsecret), nil
	}
	token, err = jwt.ParseWithClaims(t, &JwtClaimsMap{}, signingSecret)
	if err != nil {
		log.Println(err)
	}
	return token, err
}

func Parse(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	var req request.Request
	req = rc.Context().Value("request").(request.Request)
	token, err := parseToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		return err, rc
	}

	if token.Valid {
		// set scope to string "ALL" in request context if requester has privilege to all keys.
		// else, set scope to slice of uint64 key ids from the token.
		if token.Claims.(*JwtClaimsMap).Expiration < time.Now().Unix() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "text/plain")
			return errors.New("Token has expired"), rc
		}

		var mycontext context.Context

		var claims []string
		claims = token.Claims.(*JwtClaimsMap).Keys
		mycontext = context.WithValue(rc.Context(), "claims", claims)

		if token.Claims.(*JwtClaimsMap).Keys[0] == "ALL" {
			mycontext = context.WithValue(mycontext, "claimsAll", true)
		} else {
			mycontext = context.WithValue(mycontext, "claimsAll", false)
		}

		var admin int
		admin = token.Claims.(*JwtClaimsMap).Admin
		mycontext = context.WithValue(mycontext, "admin", admin)

		return nil, *rc.WithContext(mycontext)
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "text/plain")
	return errors.New("Token is not valid"), rc
}

// IsAdmin Only continue if the requestor is an admin.
func IsAdmin(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	admin := rc.Context().Value("admin").(int)
	if admin >= 1 {
		return nil, rc
	}
	w.WriteHeader(http.StatusUnauthorized)
	return errors.New("Requestor is not admin"), rc
}

// Info return jwt token information.
func Info(w http.ResponseWriter, rc http.Request) (error, http.Request) {
	type TokenInfo struct {
		Keys  []string `json:"keys"`
		Admin int      `json:"admin"`
	}
	JsonEncode := json.NewEncoder(w)

	var keys []string
	switch rc.Context().Value("claims").(type) {
	case []string:
		keys = rc.Context().Value("claims").([]string)
	case string:
		keys = append(keys, rc.Context().Value("claims").(string))
	}

	tokenInfo := TokenInfo{
		Keys:  keys,
		Admin: rc.Context().Value("admin").(int),
	}

	w.Header().Set("Content-Type", "text/json")
	JsonEncode.Encode(tokenInfo)
	return nil, rc
}

// IssueJWT return a valid jwt with these statically defined scope values.
func New(w http.ResponseWriter, rc http.Request) (error, http.Request) {

	var req request.Request
	req = rc.Context().Value("request").(request.Request)

	authKeys := rc.Context().Value("claims").([]string)
	authAllKeys := rc.Context().Value("claimsAll").(bool)

	if rc.Context().Value("admin").(int) >= 0 {
		if !authAllKeys {
			if !utils.AllStringsInSlice(req.Keys, authKeys) {
				log.Println("Requested Keys are not in requestors allowed")
				return errors.New("Permission Denied.  Requested Keys are not in requestors allowed scope"), rc
			}
		}
	} else {
		return errors.New("Permission Denied.  Not an admin token"), rc
	}

	tokenString, err := GenerateToken(req.Keys, time.Now().Unix()+req.Expiration, randid.Generate(32), req.Admin)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, rc
	}

	io.WriteString(w, tokenString)
	return nil, rc
}

func GenerateToken(keys []string, exp int64, jti string, admin int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"keys":  keys,
		"exp":   exp,
		"jti":   jti,
		"admin": admin,
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWTsecret))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, err
}
