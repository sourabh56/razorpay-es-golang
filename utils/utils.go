package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(stringLength int) string {
	var possibleCombinations = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // possible string combination
	unixTime := time.Now().UnixNano()                                                                   // rand fn behaves same after each restart so pass timeStamp to maintain uniqueness

	if stringLength <= 19 {
		return fmt.Sprint(unixTime)[0:stringLength]
	}
	source := rand.NewSource(unixTime)
	r := rand.New(source)
	stringLength = stringLength - 19                // as unixTime is of length 19 so for better safety will append 19 unixTime while returning
	randomStringArray := make([]rune, stringLength) // creating empty array of required length

	for i := 0; i < stringLength; i++ { // required length - 19(UnixNano) as appending 19 digits when returning
		randomValue := r.Intn(len(possibleCombinations))         // getting random value between 62
		randomStringArray[i] = possibleCombinations[randomValue] // getting value at that position in possibleCombinations and adding it in random string
	}

	return (string(randomStringArray) + fmt.Sprint(unixTime)) // final value
}
