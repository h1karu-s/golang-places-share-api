package controllers

import (
	// "fmt"
	"net/http"
	"time"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	jwt "github.com/dgrijalva/jwt-go"
	// "github.com/gorilla/mux"
	
	"../database"
	"../models"
	"../util"
)


//GetUsers .
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users models.Users
	cur, err := database.UserColl.Find(context.Background(), bson.D{})
	if err != nil {
		  util.Err(w, err, http.StatusInternalServerError)
		  return
	}	
	if err = cur.All(context.TODO(), &users.Users); err != nil {
			util.Err(w, err, http.StatusInternalServerError)
		  return 
	}
  json.NewEncoder(w).Encode(&users)
}


//SignUp . //検証 => 既存のユーザがいるかcheck(email) => passwordをhash化 => 保存 => jwt作成 => jwtを持ったjsonを返す
func SignUp (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var user models.UserContainPassword
	json.NewDecoder(r.Body).Decode(&user)
	email := user.Email
	//check user email
  var existingUser models.User
	err := database.UserColl.FindOne(context.Background(), bson.D{{Key: "email", Value: email}}).Decode(&existingUser)
	if err == nil {
		util.ReturnWithErrJSON(w, "this email already useing", http.StatusBadRequest)
	}
	//password hash and insert database
	h := hmac.New(sha256.New, []byte("mysecretkey"))
	h.Write([]byte(user.Password))
	hashedPassword := h.Sum(nil)
	user.Password = hex.EncodeToString(hashedPassword)
	//またこの時点でuser構造体のplacesはnilになるため[]を代入する
	user.Places = make([]string, 0)
	// var savedUser models.User
	u, err := database.UserColl.InsertOne(context.Background(), user)
	if err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return 
	}
	uid := u.InsertedID
	//create jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": uid,
		"email": email,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	//add secret key
	tokenString, err := token.SignedString([]byte("mysecretkey"))
	if err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return
	}
	res := models.SignUpUser{UserID: uid.(primitive.ObjectID), Email: user.Email, Token: tokenString}
  json.NewEncoder(w).Encode(res)
}


//Login . validation => user check for email => password check => create jwt and return json containg jwtToken 
func Login (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//validation ....
	var loginUser models.LoginUser
	json.NewDecoder(r.Body).Decode(&loginUser)
	if (loginUser.Email == "" || loginUser.Password == "" ) {
		util.ReturnWithErrJSON(w, "invalid value", http.StatusBadRequest)
    return 
	}
	//user check for email
	var user models.UserContainPassword
	result := database.UserColl.FindOne(context.Background(), bson.D{{Key: "email", Value: loginUser.Email}})
  if err := result.Decode(&user); err != nil {
		util.ReturnWithErrJSON(w, "user is not found.", http.StatusBadRequest)
		return
	}
	//password check!
	h := hmac.New(sha256.New, []byte("mysecretkey"))
	h.Write([]byte(loginUser.Password))
	hashedPassword := h.Sum(nil)
	comperedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return
	}
	if !hmac.Equal(hashedPassword, comperedPassword) {
		util.ReturnWithErrJSON(w, "password is invalid.", http.StatusBadRequest)
		return
	}

	//create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"email" : user.Email,
		"exp"   : time.Now().Add(time.Hour).Unix(),
	})
	//add secret key
	tokenString, err := token.SignedString([]byte("mysecretkey"))
	if err != nil {
		util.Err(w, err, http.StatusInternalServerError)
		return
	}
	res := models.SignUpUser{UserID: user.ID, Email: user.Email, Token: tokenString}
	json.NewEncoder(w).Encode(res)
}



