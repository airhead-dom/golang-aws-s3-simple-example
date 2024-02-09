package rest

import (
	"airhead-dom/golang-aws-s3/pkg/objectstorage"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Run() {
	clientId := viper.GetString("CLIENT_ID")
	clientSecret := viper.GetString("CLIENT_SECRET")
	bucketName := viper.GetString("BUCKET_NAME")

	sos, err := objectstorage.NewS3ObjectStorage(clientId, clientSecret, bucketName)
	if err != nil {
		log.Printf("failed loading s3 object storage. err: %v", err)
		os.Exit(1)
	}

	router := gin.Default()

	router.GET("/objects", func(ctx *gin.Context) {
		keys, err := sos.ListObject(ctx.Request.Context())

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    keys,
		})
	})

	router.GET("/single-object", func(ctx *gin.Context) {
		key := ctx.Request.URL.Query().Get("key")

		if key == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "please provide a key",
			})
			return
		}

		link, err := sos.GetObject(ctx.Request.Context(), key)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("failed getting object. err: %v", err),
			})
			return
		}

		escaped := url.QueryEscape(link)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    escaped,
		})

		return
	})

	router.Run()
}
