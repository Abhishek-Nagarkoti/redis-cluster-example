package handlers

import (
	"github.com/chasex/redis-go-cluster"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
)

type Handler struct {
	DB *redis.Cluster
}

/*=================================
***   establish connection  ***
=================================*/
func (h *Handler) Connect() {
	var err error
	h.DB, err = redis.NewCluster(
		&redis.Options{
			StartNodes: []string{os.Getenv("HOST") + ":" + os.Getenv("PORT1"), os.Getenv("HOST") + ":" + os.Getenv("PORT2"), os.Getenv("HOST") + ":" + os.Getenv("PORT3")},
		})
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
}

/*-----  End of connect  ----*/

/*=================================
***   set value in database  ***
=================================*/

func (h *Handler) Set(ctx *gin.Context) {
	body := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"Error": "Validation error"})
	} else {
		_, err = h.DB.Do("set", body.Key, body.Value)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		} else {
			ctx.JSON(200, gin.H{"key": body.Key, "value": body.Value})
		}
	}
}

/*-----  End of set  ----*/

/*===================================
***   get value from database  ***
===================================*/

func (h *Handler) Get(ctx *gin.Context) {
	value, _ := ctx.GetQuery("key")
	if value == "" {
		ctx.JSON(400, gin.H{"Error": "Wrong query string"})
	} else {
		reply, err := redis.String(h.DB.Do("get", value))
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		} else {
			ctx.JSON(200, gin.H{"key": value, "value": reply})
		}
	}
}

/*-----  End of get  ------*/
