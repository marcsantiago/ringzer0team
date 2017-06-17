package main

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func generateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)
	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan string) {
		defer close(c)
		addLetter(c, "", alphabet, length)
	}(c)

	return c
}

// AddLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
func addLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}
	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}

func main() {
	// known length of password is 6 characters
	// abcdefghijklmnopqrstuvwxyz0123456789
	for perm := range generateCombinations("abcdefghijklmnopqrstuvwxyz0123456789", 6) {
		h1 := sha1.New()
		// fmt.Println(perm)
		io.WriteString(h1, perm)
		h := fmt.Sprintf("%x", h1.Sum(nil))
		if h == "58066ab1f1ffb482737dd9162b35d9cdad5f86c8" { //spzedg
			fmt.Println(perm)
			break
		}
	}
}
