package models

import(
	"log"
)

//Handels errors
func Perror(err error) {
    if err != nil {
		log.Println(err)
        panic(err)
    }
}