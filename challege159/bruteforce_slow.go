package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type set struct {
	Hash map[string]error
	Mtx  sync.Mutex
}

var mySet set

func init() {
	mySet.Hash = make(map[string]error)
}

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
// this would create somthing like a 17.414258688 GB file and also 17.414258688 in memory... so use at your own risk
func addLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}
	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		mySet.Mtx.Lock()
		if _, ok := mySet.Hash[newCombo]; !ok {
			mySet.Hash[newCombo] = nil
			c <- newCombo
		} else {
			c <- newCombo
		}
		mySet.Mtx.Unlock()
		addLetter(c, newCombo, alphabet, length-1)
	}
}

func main() {
	for {
		select {
		case c := <-generateCombinations("abcdefghijklmnopqrstuvwxyz0123456789", 6):
			_ = c
		default:
			b, err := json.MarshalIndent(mySet.Hash, "", "")
			if err != nil {
				log.Fatalln("could not marshal data")
			}
			err = ioutil.WriteFile("hashmap.json", b, 0666)
			if err != nil {
				log.Fatalln("could not write file")
			}
		}
	}
}
