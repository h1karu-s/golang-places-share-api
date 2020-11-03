package middlewares

import (
	"net/http"
	"strings"
	"fmt"
	
	// "github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"
)

//UserID read from jwt payload.
var UserID string

//CheckAuth .
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//UserID初期化
		UserID = ""
		//urlからmatchしたものだけを適用する...
		//これ以外の場合はauthorizationをcheckする
		if r.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			w.WriteHeader(http.StatusOK)
			return 
		}
		if ((r.Method == "GET") || (r.URL.String() == "/api/users/login") || (r.URL.String() == "/api/users/signup") ) {
			fmt.Println("auth return")
			next.ServeHTTP(w, r)
			return
		} 
		if r.Header["Authorization"] == nil  {
			w.Write([]byte("Authorization header none!"))
			return
		}
		tokenAll := strings.Split(r.Header["Authorization"][0], " ") 
		if len(tokenAll) < 2 {
			w.Write([]byte("Authorization header is invalid!"))
			return
		}
		var tokenString string
		if len(tokenAll) == 2 {
			  tokenString = tokenAll[1]	
		} else {
				fmt.Println("tokensize is not invalid.")
				return
		}
		if tokenString == " " {
			w.Write([]byte("token not found!!"))
			return
		}
		//tokenCheck
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("mysecretkey"), nil
		})
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(token)
			w.Write([]byte(err.Error()))
			return
		}
		if token == nil {
			w.Write([]byte("invalid token! nil"))
			return
		}
		if !token.Valid || err != nil {
			w.Write([]byte("invalid token!"))
			return
		}
		if (token.Valid) {
			UserID = token.Claims.(jwt.MapClaims)["userId"].(string)
		}

		fmt.Println(tokenString)
			next.ServeHTTP(w, r)
	})
}
