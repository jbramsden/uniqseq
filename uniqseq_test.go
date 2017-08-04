//uniqseq_test.go
package uniqseq

import (
	"testing"
	"time"
)

func TestCreateAndInit(t *testing.T) {
	a := Create()
	a.Init()
	if a.counter != -1 {
		t.Error("Init should have set the counter to -1")
	}
	if len(a.jumblerSet) != 0 {
		t.Error("Jumble character set and a string when the default value should not set this")
	}
	x := a.Encode(918756211)
	t.Errorf("%s", x)
	a.Jumbler = true
	a.Init()
	if len(a.jumblerSet) != len(a.CharacterSet) {
		t.Error("Jumble set and character set should be the same length")
	}

}

func TestEncode(t *testing.T) {
	a := Create()
	a.CharacterSet = "ABCDEF"
	a.BlankFillChar = ""
	if a.Encode(5) != "F" {
		t.Error("Expecting F but got ", a.Encode(5))
	}
	a.BlankFillChar = "A"
	if a.Encode(0) != "AAAA" {
		t.Error("Expecting AAAA but got ", a.Encode(0))
	}
	if a.Encode(6) != "AABA" {
		t.Error("Expecting AA but got ", a.Encode(6))
	}

}

func TestNext(t *testing.T) {
	a := Create()
	a.CharacterSet = "abcd"
	a.BlankFillChar = ""
	a.Init()
	c, d := a.Next()
	if c != "a" || d != 0 {
		t.Error("Expect the first character in the sequence 'a' with a number of 0 and got ", c, d)
	}
	c, d = a.Next()
	if c != "b" || d != 1 {
		t.Error("Expect the first character in the sequence 'b' with a number of 1 and got ", c, d)
	}
	c, d = a.Next()
	if c != "c" || d != 2 {
		t.Error("Expect the first character in the sequence 'c' with a number of 2 and got ", c, d)
	}
	c, d = a.Next()
	if c != "d" || d != 3 {
		t.Error("Expect the first character in the sequence 'd' with a number of 3 and got ", c, d)
	}
	c, d = a.Next()
	if c != "ba" || d != 4 {
		t.Error("Expect the first character in the sequence 'ba' with a number of 4 and got ", c, d)
	}

	//testing concurreny by starting 10 go subroutines
	a.Init()
	for i := 0; i < 10; i++ {
		go func() {
			for x := 0; x < 5; x++ {
				a.Next()
			}
		}()
	}
	//wait for all go routines to finish so 5 seconds should be more then enough.
	time.Sleep(5 * time.Second)
	//10 go subroutines * 5 calls for each to Next() means 50, but minus 1 as the counter starts from 0
	if a.counter != 49 {
		t.Error("Expecting 49 on the counter but got", a.counter)
	}
}
