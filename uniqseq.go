// Licensed under the MIT License(the "License");

// Package uniqseq is a simple package to generate a unique string that can be used for URL shortners or for anything else where you need to have a ordered unique string.
// Based on the package: https://github.com/xor-gate/bjf
package uniqseq

import (
	"strings"
	"sync"
	"time"
)

//Base is to handle the characters to be used within the unique string
type Base string

const (
	//Full is for the full english character set
	Full Base = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//NoVowels is the full english character set minus vowles. This should be used so that no rude words appear!
	NoVowels Base = "bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ123456789"
)

//UniqueString - holds the main settings and functions for generating a unique string
type UniqueString struct {
	CharacterSet  Base  //This is the character set to use when generating a unique string
	StartLength   int   //The starting length of the unique string
	BlankFillChar Base  //This is the single character that will be used to fill the blank spaces when the generated string length does not match the StartLength
	Jumbler       bool  //If set this will jumble up the CharacterSet on Initialization and the JumblerSet will be used instead of the CharacterSet
	jumblerSet    Base  //The same as CharacterSet except this has been Jumbled up if the Jumbler flag is set to true before Initalization
	LastCharInc   bool  //If set to true the last character will increment[aaaa,aaab,aaac....]. If false the first character will increement[aaaa,baaa,caaa,...]
	counter       int64 //holds the counter for each next sequence number
	countLock     sync.RWMutex
	LoopAt        int64 //If set to anything greater then 0, when the counter reaches this value it will start the counter again at 1.
}

//Create - returns a uniqueString structure with default values
func Create() *UniqueString {
	var def UniqueString
	def.CharacterSet = Full
	def.StartLength = 4
	def.BlankFillChar = "_"
	def.Jumbler = false
	def.LastCharInc = true
	return &def
}

//Init - Initalisation of the counter and creates the jumbler of the character set if required.
func (u *UniqueString) Init() {
	u.counter = -1
	if u.Jumbler {
		u.jumble()
	}
}

//Next - Return the next unique string and coresponding number in the sequence
func (u *UniqueString) Next() (string, int64) {
	u.countLock.Lock()
	u.counter++
	if u.LoopAt > 0 && u.counter >= u.LoopAt {
		u.counter = 0
	}
	value := u.Encode(u.counter)
	u.countLock.Unlock()
	return value, u.counter
}

//jumble - This will jumble up the Character set used. This is to minimise anyone guessing the next sequence of the unique string
func (u *UniqueString) jumble() {
	td := time.Now().UnixNano()
	//New character set
	ncs := make([]byte, 0)
	//Working character set
	wcs := string(u.CharacterSet)

	//decide first character
	p := td % int64(len(wcs))
	ncs = append(ncs, wcs[p])
	wcs = strings.Replace(wcs, string(wcs[p]), "", -1)

	for i := 1; i < len(u.CharacterSet); i++ {
		n := (p * int64(i)) % int64(len(wcs))
		ncs = append(ncs, wcs[n])
		wcs = strings.Replace(wcs, string(wcs[n]), "", -1)
	}
	u.jumblerSet = Base(ncs)
}

//Encode - This converts the passed in interger to a unique string.
func (u *UniqueString) Encode(n int64) string {

	t := make([]byte, 0)

	cSet := u.getCharSet()

	if n == 0 {
		t = append(t, cSet[0])
	}

	lenaC := int64(len(cSet))
	for n > 0 {

		r := n % lenaC
		t = append(t, cSet[r])
		n = n / lenaC
	}

	if len(u.BlankFillChar) > 0 {
		for len(t) < u.StartLength {
			t = append(t, u.BlankFillChar[0])
		}
	}

	if u.LastCharInc {
		for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
			t[i], t[j] = t[j], t[i]
		}
	}

	return string(t)
}

//getCharSet - returns the character set to use.
func (u *UniqueString) getCharSet() Base {
	ac := u.CharacterSet
	if u.Jumbler && len(u.jumblerSet) > 0 {
		ac = u.jumblerSet
	}
	return ac
}
