package web

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/aja-video/contra/src/configuration"
	"github.com/go-ldap/ldap"
	"log"
	"net/http"
)

// auth holds authentication methods
type auth interface {
	authenticate(nextHandler http.HandlerFunc) http.HandlerFunc
	// logout TODO
}

// basicAuth implements auth with http basic authentication
type basicAuth struct {
	username string
	password string
}

// ldap implements auth with ldap authentication
type ldapAuth struct {
	server   string
	username string
	password string
	filter   string
	search   ldap.SearchRequest
}

// buildAuth initializes auth based on the passed configuration
func (w *Web) buildAuth(c *configuration.Config) auth {
	switch c.WebAuth {
	case "basic":
		return &basicAuth{
			username: c.WebUser,
			password: c.WebPass,
		}
	case "ldap":
		return &ldapAuth{
			server:   c.LdapServer,
			username: c.LdapBindUser,
			password: c.LdapBindPass,
			filter:   c.LdapFilter,
			search: ldap.SearchRequest{
				BaseDN:       c.LdapBaseDN,
				Scope:        2,
				DerefAliases: 0,
				SizeLimit:    0,
				TimeLimit:    0,
				TypesOnly:    false,
				Filter:       "", // We'll populate this on the request when we have a username
				Attributes:   []string{},
				Controls:     nil,
			},
		}
	default:
		return nil
	}
}

// authenticate (basic) handles http basic authentication
func (b *basicAuth) authenticate(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(b.username))
			expectedPasswordHash := sha256.Sum256([]byte(b.password))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				nextHandler.ServeHTTP(res, req)
				return
			}
		}

		res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
	}
}

// authenticate (ldap) handles ldap authentication
func (b *ldapAuth) authenticate(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		l, err := ldap.DialURL(b.server)
		if err != nil {
			log.Println(err)
			res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		err = l.Bind(b.username, b.password)
		if err != nil {
			log.Println(err)
			res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		defer l.Close()
		b.search.Filter = fmt.Sprintf("(&%s(SamAccountName=%s))", b.filter, user)
		if ok {
			search, err := l.Search(&b.search)
			if err != nil {
				log.Println(err)
				res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
				http.Error(res, err.Error(), http.StatusUnauthorized)
				return
			}
			if len(search.Entries) != 1 {
				fmt.Println(len(search.Entries))
				res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
				http.Error(res, "Incorrect Username or Password", http.StatusUnauthorized)
				return
			}
			//fmt.Println(search.Entries[0].DN)
			if err = l.Bind(search.Entries[0].DN, pass); err != nil {
				log.Println(err)
				res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
			nextHandler.ServeHTTP(res, req)
			return
		}
		res.Header().Set("WWW-Authenticate", `Basic realm="CONTRA", charset="UTF-8"`)
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
	}
}
