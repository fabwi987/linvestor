package controllers

import "log"

//Perror handels errors
func Perror(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
