package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Sample from: https://go.dev/doc/tutorial/fuzz

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev := Reverse(input)
	doubleRev := Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}

func Reverse(s string) string {
	return ReverseV2(s)
}

/*
Flaw: Characters such as 泃,ĝ can require several bytes. Thus, reversing the string byte-by-byte will invalidate multi-byte characters.

go test -fuzz=Fuzz
--- FAIL: FuzzReverse (0.00s)

	reverse_test.go:32: Number of runes: orig=1, rev=2, doubleRev=1
	reverse_test.go:37: Reverse produced invalid UTF-8 string "\x9d\xc4"
*/
func ReverseV1(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

/*
Fixed from V1: The key difference is that Reverse is now iterating over each rune in the string, rather than each byte.
Flaw: Output the string is different from the original after being reversed twice. This time the input itself is invalid unicode.

go test -run=FuzzReverse/49d56c4906d9c4bd
--- FAIL: FuzzReverse (0.00s)

	reverse_test.go:32: Number of runes: orig=1, rev=1, doubleRev=1
	reverse_test.go:34: Before: "\x9f", after: "�"
*/
func ReverseV2(s string) string {
	fmt.Printf("input: %q\n", s)
	r := []rune(s)
	fmt.Printf("runes: %q\n", r)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

/*
Fixed from V2: Reverse is now returning an error if the input is not valid UTF-8

go test -run=FuzzReverseV3  -fuzz=FuzzReverseV3 -fuzztime 30s
fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 861260 (287032/sec), new interesting: 38 (total: 41)
fuzz: elapsed: 6s, execs: 1982417 (373745/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 9s, execs: 3026971 (348175/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 12s, execs: 4026657 (333267/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 15s, execs: 5094136 (355785/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 18s, execs: 6103937 (336621/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 21s, execs: 7175574 (357227/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 24s, execs: 8241507 (355291/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 27s, execs: 9289180 (349171/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 30s, execs: 10311269 (340773/sec), new interesting: 42 (total: 45)
fuzz: elapsed: 30s, execs: 10311269 (0/sec), new interesting: 42 (total: 45)
PASS
ok      example/fuzz    30.243s
*/
func ReverseV3(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}
