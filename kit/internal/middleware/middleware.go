package middleware

import "github.com/go-kit/kit/log"

type Middleware func()

func LogMiddleware(logger log.Logger) {

}
