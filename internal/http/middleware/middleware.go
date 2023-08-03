package middleware

import (
	"net/http"
	"regexp"
)

func CheckMethodAndContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		// Test was crashed
		// if contentType := r.Header.Get("content-type"); contentType != "text/plain" {
		// 	http.Error(w, "Only Content-Type is text/plain!", http.StatusMethodNotAllowed)
		// 	return
		// }

		if re := regexp.MustCompile(`/update/[a-zA-Z]+/[a-zA-Z]+/\d+[/]{0,}$`); !re.MatchString(r.URL.Path) {
			http.Error(w, "Invalid URL format!", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
