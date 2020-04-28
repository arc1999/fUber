package Auth

import (
	"fUber/errs"
	"encoding/json"
	"errors"
	"fUber/Request"
	"fUber/Response"
	"fUber/types"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	APP_KEY = "DAEMONS"
)
type User_temp struct{
	Uname string `json:"Username"`
	Upass string `json:"Password"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	var u1 User_temp
	var temp types.User
	_ = json.NewDecoder(r.Body).Decode(&u1)
	fmt.Println(u1)
	db, err := gorm.Open("mysql", "root:root@tcp(mysqldb:3306)/fUber?charset=utf8&parseTime=True")
	Db:= db.Table("users").Where("u_name = ?", u1.Uname).Find(&temp)
	if Db.RowsAffected == 0{
		err=errors.New("Incorrect Credentials")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}
	if  u1.Upass != temp.U_pass {
		w.WriteHeader(http.StatusUnauthorized)
		err=errors.New("Incorrect Credentials")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}

	// an expiry of 1 hour.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": temp.U_id,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

func AuthMiddleware(next http.Handler) http.Handler {
	if len(APP_KEY) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(APP_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}
func Signup(w http.ResponseWriter, r *http.Request){
	var u1 Request.CreateUserRequest
	var u2 types.User
	_ = render.Bind(r, &u1)
	db, err := gorm.Open("mysql", "root:root@tcp(mysqldb:3306)/fUber?charset=utf8&parseTime=True")
	if(err!=nil){
		err=errors.New("Database error")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}
	defer db.Close()
	temp:=db.Table("users").Where("u_name = ?",u1.U_name).Find(&u2)
	if(temp.RowsAffected!=0){
		err=errors.New("Username Already exsist ")
		_ = render.Render(w, r, errs.ErrRender(err))
		return
	}
	u1.Amount=5000.00
	db.Create(&u1.User)
	render.Render(w,r,Response.ListUser(&u1))
}
