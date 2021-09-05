package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

type state int

const (
	begin state = iota
	escape
	symbol
)

func Unpack(s string) (string, error) {
	var prevRune rune
	state := begin

	b := strings.Builder{}

	for _, r := range s {
		switch {
		case '0' <= r && r <= '9':
			switch state {
			case begin:
				return "", ErrInvalidString
			case escape:
				prevRune = r
				state = symbol
			case symbol:
				for i := rune(0); i < r-'0'; i++ {
					b.WriteRune(prevRune)
				}
				state = begin
			}
		case r == '\\':
			switch state {
			case begin:
				state = escape
			case escape:
				prevRune = r
				state = symbol
			case symbol:
				b.WriteRune(prevRune)
				state = escape
			}
		default:
			switch state {
			case begin:
				prevRune = r
				state = symbol
			case escape:
				return "", ErrInvalidString
			case symbol:
				b.WriteRune(prevRune)
				prevRune = r
			}
		}
	}

	switch state {
	case begin:
	case escape:
		return "", ErrInvalidString
	case symbol:
		b.WriteRune(prevRune)
	}

	return b.String(), nil
}
