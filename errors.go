package yt

import "fmt"

type ErrNotImplemented struct {
	Functionality string
}

func (ni ErrNotImplemented) Error() string {
	return fmt.Sprintf("not implemented: %s", ni.Functionality)
}

type ErrEmptyToken struct{}

func (et ErrEmptyToken) Error() string {
	return fmt.Sprintf("empty token")
}

type ErrInvalidEscapeChar struct {
	Rune rune
}

func (ie ErrInvalidEscapeChar) Error() string {
	return fmt.Sprintf("invalid escape character: %v", ie.Rune)
}

type ErrUnexpectedRune struct {
	Rune rune
}

func (uc ErrUnexpectedRune) Error() string {
	return fmt.Sprintf("unexpected rune: %s", string(uc.Rune))
}

type ErrUnknown struct {
	Message string
}

func (u ErrUnknown) Error() string {
	return u.Message
}

type ErrExpectedMap struct{}

func (em ErrExpectedMap) Error() string {
	return "expected map"
}

type ErrKeyNotFound struct {
	Key string
}

func (knf ErrKeyNotFound) Error() string {
	return fmt.Sprintf("not found: %s", knf.Key)
}

type ErrExpectedArray struct{}

func (ea ErrExpectedArray) Error() string {
	return "expected array"
}

type ErrOutOfBounds struct{}

func (oob ErrOutOfBounds) Error() string {
	return "index out of bounds"
}

type ErrKeyAlreadyDefined struct {
	Key string
}

func (kc ErrKeyAlreadyDefined) Error() string {
	return fmt.Sprintf("key already defined: %s", kc.Key)
}

type ErrInvalidQuery struct{}

func (iq ErrInvalidQuery) Error() string {
	return "invalid query"
}

type ErrCycleDetected struct {
	Source string
}

func (cd ErrCycleDetected) Error() string {
	return fmt.Sprintf("cycle detected: %s", cd.Source)
}
