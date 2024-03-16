package helper

import (
	//"go/token"
	"os"

	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)
func init(){
	godotenv.Load()
	
}
var SECRET_KEY string = os.Getenv("SECRET_KEY")
var REFRESH_SECRET_KEY string = os.Getenv("REFRESH_SECRET_KEY")
type signedDetail struct {
	Email     string
	Name      string
	User_type string
	jwt.RegisteredClaims
}
func GenerateAlltoken(email string, name string, user_type string, db *gorm.DB)(signedToken string, signedRefreshToken string, err error)  {
	var SECRET_KEY string = os.Getenv("SECRET_KEY")
	var REFRESH_SECRET_KEY string = os.Getenv("REFRESH_SECRET_KEY")
	claims := &signedDetail{
		Email: email,
		Name: name,
		User_type: user_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24)),
		} ,
	}
	refreshClaims:= &signedDetail{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 168)),
		},
	}
	token,err:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	if err!=nil {
		log.Panic(err)
		return 
	}
	refreshToken, err:= jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SignedString([]byte(REFRESH_SECRET_KEY))
	if err!=nil {
		log.Panic(err)
		return 
	}
	return token,refreshToken,err
}

func IsAuthorized(tokenString string) (*signedDetail, error) {
	claims := &signedDetail{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil // Use your secret key here
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if the token is not expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}
