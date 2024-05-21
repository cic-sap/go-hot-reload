package main

import _ "embed"
import (
	"fmt"
	reload "github.com/cic-sap/go-hot-reload/v1"

	"log"
	"net/http"
	"os"
	"time"
)

import "github.com/gin-gonic/gin"

//go:embed ver.txt
var ver string

func main() {

	log.Println("start")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	go reload.AutoReload()

	// your app code ...
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("ver", ver, "pid", os.Getpid())
		}
	}()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("hello ver:%s", ver))
	})

	r.Run("0.0.0.0:2345")

	log.Println("end")
}
