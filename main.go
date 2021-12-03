package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type quiry struct {
	Title       string
	Author      string
	Date        string
	Description string
	Page        int
	Sha256      string
	Searchfor   string
}

func Latest(c *gin.Context) {
	var student quiry
	// 读取请求数据并写入结构体b
	c.Bind(&student)
	// 返回 JSON 格式响应
	c.JSON(200, gin.H{
		"Title":       student.Title,
		"Author":      student.Author,
		"Date":        student.Date,
		"Description": student.Description,
		"Page":        student.Page,
	})
}
func Page(c *gin.Context) {
	var student quiry
	// 读取请求数据并写入结构体b
	c.Bind(&student)
	// 返回 JSON 格式响应
	data := gin.H{"Sha256": student.Sha256}
	c.JSON(200, data)
}
func Content(c *gin.Context) {
	var student quiry
	// 读取请求数据并写入结构体b
	c.Bind(&student)
	// 返回 JSON 格式响应
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")	//要改
	if uri == "" {		//要改
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("offices").Collection("dean")
	title := "Back to the Future"	//要改
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"sha256", student.Sha256}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)	//要改
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{"Data": jsonData})
}
func main() {
	r := gin.Default()
	r.GET("/latest", Latest)
	r.GET("/page", Page)
	r.GET("/text", Content)
	r.Run()
}
