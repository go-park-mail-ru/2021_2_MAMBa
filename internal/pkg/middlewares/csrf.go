package middlewares

import "github.com/gorilla/csrf"

var authKey = []byte("32-byte-long-auth-key")

var CSRF = csrf.Protect(
	authKey,
	csrf.Path("/"),
	csrf.Secure(false)) // TODO: Сделать Secure в продакшн
