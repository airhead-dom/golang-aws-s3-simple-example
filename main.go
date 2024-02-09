package main

import (
	"airhead-dom/golang-aws-s3/rest"
	"log"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("local")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("failed reading config %v", err)
		os.Exit(1)
	}

	rest.Run()
}
