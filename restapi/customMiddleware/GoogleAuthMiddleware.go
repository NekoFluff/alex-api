package customMiddleware

import (
	"addi/utils"
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/idtoken"
)

func GoogleAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")

		authTokens := r.Header.Values("google-idToken")
		if len(authTokens) > 0 {
			token := authTokens[0]
			googleClientId := utils.GetEnvVar("GOOGLE_CLIENT_ID")

			payload, err := idtoken.Validate(context.Background(), token, googleClientId)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			fmt.Println(payload)
			ctx := context.WithValue(r.Context(), "google-idToken", payload)

			// email := payload.Claims["email"]
			// name := payload.Claims["name"]

			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			w.Write([]byte("Missing google-idToken header. Unauthorized request."))
			return
		}

	})
}
