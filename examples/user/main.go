package main

import (
	"fmt"
	"log"

	"github.com/avocatl/admiral/pkg/prompter"
)

type User struct {
	FirstName          string
	LastName           string
	Age                int `prompter:"-"`
	Username           string
	TermsAndConditions bool
	Newsletter         bool `prompter:"-"`
}

func main() {
	u, err := prompter.Struct(&User{FirstName: "John", Age: 17})
	if err != nil {
		log.Fatal(err)
	}

	msg, err := prompter.String("Additional messages", "")
	if err != nil {
		log.Fatal(err)
	}

	usr := u.(*User)

	fmt.Printf(`
Your user %s aged %d left a message, 
to reply tag him with @%s.
The message is:
%s
`, usr.FirstName, usr.Age, usr.Username, msg)
}
