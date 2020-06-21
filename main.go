package main

import (
	"github.com/hatzelencio/branch-protection/remote"
	"log"
)

func main() {
	err := remote.ValidateInputs()

	if err != nil {
		log.Fatal(err)
	}

	err = remote.UpdateBranchProtection()
	if err != nil {
		log.Fatal(err)
	}
}
