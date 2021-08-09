package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"strconv"
	"time"
)

const SecretKey = "secret"

func CheckPasswordHash(password []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//make job
func Todo(c *gin.Context) {

	// var json struct {
	// 	tokenid string `json:"tokenid" `
	// }
	requestBody, _ := ioutil.ReadAll(c.Request.Body)
	bodyJson := make(map[string]interface{}) //注意该结构接受的内容
	json.Unmarshal(requestBody, &bodyJson)
	fmt.Println("验证通过", bodyJson["tokenid"])
	str := fmt.Sprintf("%v", bodyJson["tokenid"])
	c.String(200, "验证通过"+str)
}

//注册
func Register(c *gin.Context) {

	requestBody, err := ioutil.ReadAll(c.Request.Body)

	var data map[string]interface{}
	err = json.Unmarshal(requestBody, &data)
	var account entity.Account

	if err != nil {
		c.JSON(http.StatusBadRequest, "参数异常")
		return
	}

	password := []byte(data["Password"].(string))
	hash, _ := bcrypt.GenerateFromPassword(password, 14)
	account.Password = hash
	account.Name = data["Name"].(string)
	account.Email = data["Email"].(string)

	database.Connector.Create(&account)
	c.JSON(http.StatusOK, "完成创建")

}

//登录
func Login(c *gin.Context) {

	var data map[string]interface{}

	requestBody, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(requestBody, &data)
	fmt.Println(data)
	var accountDB entity.Account
	database.Connector.Where("email = ?", data["Email"]).First(&accountDB)

	if accountDB.Id == 0 {
		c.JSON(200, "账号不存在")
		return
	}

	password := []byte(data["Password"].(string))

	if CheckPasswordHash(password, accountDB.Password) != true {
		c.JSON(200, "密码错误")
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(accountDB.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		fmt.Println("err != nil")
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		HttpOnly: true,
	}

	fmt.Println("Cookie  create", token)

	//http.SetCookie(w, &cookie)
	http.SetCookie(c.Writer, &cookie)
	fmt.Println("cookie", cookie)
	fmt.Println("--------------------------")
	m3 := map[string]string{
		"message": token,
		"result":  "ok",
	}

	c.JSON(http.StatusOK, m3)
}

//退出
func Logout(c *gin.Context) {

	cookie, _ := c.Cookie("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	m3 := map[string]string{
		"message": "success",
	}

	if err != nil {
		fmt.Println(err)
		m3["message"] = "StatusUnauthorized"
		c.JSON(http.StatusUnauthorized, m3)
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)

	m3["Issuer"] = claims.Issuer

	newcookie := http.Cookie{
		Name:     "jwt",
		Expires:  time.Now().Add(-time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, &newcookie)
	c.JSON(200, m3)

}
