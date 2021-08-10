package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	//"go-charset"
	//"golang.org/x/text/encoding/simplifiedchinese"
	//"golang.org/x/text/transform"
	//	"github.com/gin-gonic/gin/binding"
	"bytes"
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"rest-go-demo/controllers"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"time"
)

func main() {

	initDB()
	log.Println(" 中文 Starting the HTTP server on port 8090")

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://127.0.0.1:8080", "http://localhost:8080"},
	// 	AllowCredentials: true,
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: false,
	// })

	// Insert the middleware

	r := gin.Default()

	r.Use(CORSMiddleware())

	initaliseHandlers(r)

	r.LoadHTMLGlob("views/*")
	r.GET("/views/select_file.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_file.html", gin.H{})
	})

	r.Run(":8090")
}

func costTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()

		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("the request URL %s cost %v\n", url, costTime)
	}
}

func verification() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, _ := c.Cookie("jwt")
		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {

			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, "StatusUnauthorized1")
			c.Abort()
		}

		claims := token.Claims.(*jwt.StandardClaims)
		Issuer := claims.Issuer

		// 解析并验证 JSON 格式请求数据

		if Issuer == "" {
			fmt.Println("Issuer.Length == 0")
			c.JSON(http.StatusUnauthorized, "StatusUnauthorized11")
			c.Abort()
		}
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		bodyJson := make(map[string]interface{}) //注意该结构接受的内容
		//c.ShouldBindBodyWith(&bodyJson, binding.JSON)

		json.Unmarshal(requestBody, &bodyJson)

		fmt.Println("c.Request.Header[ tokenid =", bodyJson)
		fmt.Println("Issuer =", Issuer)
		if Issuer != bodyJson["tokenid"] {
			c.JSON(http.StatusUnauthorized, "StatusUnauthorized2")
			c.Abort()
		}

		//请求处理
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	basePath := "./upload/"
	filename := basePath + filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("文件 %s 上传成功 ", file.Filename))

}

func initaliseHandlers(router *gin.Engine) {
	router.GET("/get", controllers.GetAllPerson)
	router.POST("/create", controllers.CreatePerson)
	router.GET("/get/:name", controllers.GetPersonByID)
	router.PUT("/update/:name", controllers.UpdatePersonByID)
	router.DELETE("/delete/:name", controllers.DeletPersonByID)

	router.POST("/login", controllers.Login)
	router.GET("/logout", costTime(), controllers.Logout)
	router.POST("/register", controllers.Register)
	router.GET("/todo", verification(), controllers.Todo)
	router.POST("/upload", Upload)

	router.GET("/getwater_d/:name", controllers.GetWaterDByID)
	router.GET("/water_level/:name", controllers.GetWaterLevelByID)
	router.GET("/water_ph/:name", controllers.GetWaterLevelByID)
	router.GET("/water_FluorideIon/:name", controllers.GetWaterLevelByID)
	router.GET("/water_DissolvedSolids/:name", controllers.GetWaterLevelByID)
	router.GET("/water_salinity/:name", controllers.GetWaterLevelByID)
	router.GET("/water_node", controllers.GetWaterNode)

}

func initDB() {
	config :=
		database.Config{
			ServerName: "localhost:3306",
			User:       "root",
			Password:   "qazWSX1!",
			DB:         "infodb", //infodb
		}

	connectionString := database.GetConnectionString(config)

	err := database.Connect(connectionString)
	if err != nil {
		log.Println("err != nil")
		panic(err.Error())
	}

	database.Migrate(&entity.Person{})
	database.Migrat2(&entity.Account{})
	database.Migrat3(&entity.WaterD{})
}
