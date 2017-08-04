# uniqseq

[![GoDoc](https://godoc.org/github.com/jbramsden/uniqseq?status.svg)](https://godoc.org/github.com/jbramsden/uniqseq)

This Golang package is used to create a unique sequence of characters that is usable by humans, e.g in a URL Shortner.
The theory is that the sequence ["balaCJ"] is easier to remember and type into a device then 918756211 !

The code used here is based on the Golang code of: https://github.com/xor-gate/bjf which is based on the algroithm from
http://stackoverflow.com/questions/742013/how-to-code-a-url-shortener

Unlike with xor-gate bjf package I have not included a way to reverse the unique sequence of character to get a integer, this is only because I did not need one, as I will be storing the unique string generate in a database.

### Installation

		go get github.com/jbramsden/uniqseq
		
### Quick start

The following examples give you the simplest way to use this package:
		
		package main 
		
		import (
			"fmt"
			"uniqseq"
		)
		
		func main() {
			//Create returns the uniqseq structure with default values
			us := uniqseq.Create()
			//Init resets the counter value in this simple instance
			us.Init()
			//Next get the next sequence of the unique sequence 
			UniqueString, Numb := us.Next()
			//This should output ___a, 0
			fmt.Printf("%s, %d", UniqueString, Numb)
			
		} 

This quick example uses the standard character set [a-zA-Z0-9] and sets the the minimum characters to be 4 with a filler character of '_'.

### Features

 1. Use your own character sets for the sequence. 

	There are 2 characters sets defined in this package: uniqseq.Full (default) and uniqseq.NoVowels. 
	The Full character set has [a-zA-Z0-9] where as the NoVowels is the same but without any vowels. NoVowels should be used if you want to minize rude words appearing.
	
	Example of how to use NoVowels:
	
		us := uniqseq.Create()
		us.CharacterSet = uniqseq.NoVowels
		//This will output 
		fmt.Print(us.Encode(75))
	
	Of course you can provide your own sequence of characters to use.
	
		us := uniqseq.Create()
		us.CharacterSet = "NPMQFGJ%£!%^&!@"
		//This will output ___%
		fmt.Print(us.Encode(7))
		
 1. Set the length of characters to start with.

    This will allow you to prefix the string to a certain length. E.G instead of just 'b' being returned for the number 1 you can prefix this as a minimum of 6 characters so the output would be _____b. 
	This is to be used in conjunction with BlankFillChar.
	
		us := uniqseq.Create()
		us.StartLength = 6
		//This will output _____b
		fmt.Print(us.Encode(1))
	
 1. Change the prefix character.
  
    Continuing on from the example above, this will set the prefix character. The default is '_'
	
		us := uniqseq.Create()
		us.BlankFullChar = "A"
		//This will output AAAb
		fmt.Print(us.Encode(1))
		
 1. Jumble up the character set so that it is harder for someone to guess the next value in the sequence.

	With the standard character set it is very easy to guess the number unique sequence. If you saw 'aaab' you can easily guess the next one to be 'aaac' and 'aaad' and so on. 
	To make it harder to guess the next sequence a jumbler has been added. This will take the defined character sequence and jumble up all the characters. Each time Init is called it will create a new jumbled up character set. 
	
		us := uniqseq.Create()
		us.Jumbler = true
		us.Init()
		//First should be ___[a-zA-Z0-9] 
		first, _ := us.Next()
		//Second should again be ___[a-zA-Z0-9] but it should not be in sequence with the first.
		second, _ := us.Next()
		
 1. Change the sequence to be right to left or left to right.

	Instead of the sequence alway incrementing over the character set on the right handside, this can be switched to increament on the lefthand side.
	
		us := uniqseq.Create()
		us.LastCharInc = false
		us.Init()
		//First should be a____
		first, _ := us.Next()
		//Second should be b____
		second, _ := us.Next()
		
		LoopAt
 1. Limit the sequence so that it loops round and start again.

	This stops the string getting too long as allows the reuse of the unique string.
	
		us := uniqseq.Create()
		us.LoopAt = 5
		us.BlankFullChar = ""
		us.Init()
		//Should print abcdefabcde
		for i:=0; i<10; i++ {
			nextSeq, _ := us.Next()
			fmt.Print(nextSeq)
		}

