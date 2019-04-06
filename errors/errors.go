package errors

import "fmt"

type NotImplemented struct {
	Functionality string
}
func (ni NotImplemented) Error() string {
	return fmt.Sprintf("not implemented: %s", ni.Functionality)
}

type EmptyToken struct {}
func (et EmptyToken) Error() string {
	return fmt.Sprintf("empty token")
}

type InvalidEscapeChar struct {
	Rune rune
}
func (ie InvalidEscapeChar) Error() string {
	return fmt.Sprintf("invalid escape character %v", ie.Rune)
}

type UnexpectedChar struct {
	Char rune
}
func (uc UnexpectedChar) Error() string {
	return fmt.Sprintf("unexpected character %v", uc.Char)
}

type Unknown struct {
	Message string
}
func (u Unknown) Error() string {
	return u.Message
}

type ExpectedMap struct{}
func (em ExpectedMap) Error() string {
	return "expected map"
}

type KeyNotFound struct{
	Key string
}
func (knf KeyNotFound) Error() string {
	return fmt.Sprintf("not found: %s", knf.Key)
}

type ExpectedArray struct{}
func (ea ExpectedArray) Error() string {
	return "expected array"
}

type OutOfBounds struct{}
func (oob OutOfBounds) Error() string {
	return "index out of bounds"
}

type KeyAlreadyDefined struct{
	Key string
}
func (kc KeyAlreadyDefined) Error() string {
	return fmt.Sprintf("key already defined: %s", kc.Key)
}
