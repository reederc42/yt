package yt

import "fmt"

type NotImplementedError struct {
	Functionality string
}
func (ni NotImplementedError) Error() string {
	return fmt.Sprintf("not implemented: %s", ni.Functionality)
}

type EmptyTokenError struct {}
func (et EmptyTokenError) Error() string {
	return fmt.Sprintf("empty token")
}

type InvalidEscapeCharError struct {
	Rune rune
}
func (ie InvalidEscapeCharError) Error() string {
	return fmt.Sprintf("invalid escape character %v", ie.Rune)
}

type UnexpectedCharError struct {
	Char rune
}
func (uc UnexpectedCharError) Error() string {
	return fmt.Sprintf("unexpected character %v", uc.Char)
}

type UnknownError struct {
	Message string
}
func (u UnknownError) Error() string {
	return u.Message
}

type ExpectedMapError struct{}
func (em ExpectedMapError) Error() string {
	return "expected map"
}

type KeyNotFoundError struct{
	Key string
}
func (knf KeyNotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", knf.Key)
}
