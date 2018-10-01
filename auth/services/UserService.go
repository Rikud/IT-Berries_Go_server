package services

import (
	"IT-Berries/auth/DA"
	"IT-Berries/auth/entities"
)

func FindUserByEmail(email string) *entities.User {
	searchResult := DA.FindUserByEmail(email)
	if len(searchResult) > 1 {
		panic("More than one user with the same email.")
	}
	return searchResult[0]
}
