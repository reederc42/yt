package yt

import (
	"unicode"
)

const (
	kindKey int = iota
)

type state int
const (
	invalid state = iota
	start
	key
)

type queryElement struct {
	kind int
	key string
}

func Query(v interface{}, query string) (interface{}, error) {
	elements, err := lex(query)
	if err != nil {
		return nil, err
	}
	value, err := execQuery(v, elements)
	return value, err
}

func lex(query string) ([]queryElement, error) {
	s := invalid
	var token string
	qe := make([]queryElement, 0)
	for _, r := range query {
		switch {
		case r == '.':
			switch s {
			case invalid:
				s = start
			case start:
				return nil, EmptyTokenError{}
			case key:
				qe = append(qe, queryElement{
					kind: kindKey,
					key: token,
				})
				s = start
				token = ""
			}
		case unicode.IsNumber(r):
			switch s {
			case invalid:
				return nil, NotImplementedError{
					Functionality: "indices",
				}
			case start:
				return nil, NotImplementedError{
					Functionality: "indices",
				}
			case key:
				token += string(r)
			}
		default:
			switch s {
			case invalid:
				return nil, UnexpectedCharError{
					Char: r,
				}
			case start:
				token = string(r)
				s = key
			case key:
				token += string(r)
			}
		}
	}

	switch s {
	case invalid:
		return nil, UnknownError{
			Message: "invalid query",
		}
	case key:
		qe = append(qe, queryElement{
			kind: kindKey,
			key: token,
		})
	}
	return qe, nil
}

//parsing is implied by query elements
func execQuery(v interface{}, query []queryElement) (interface{}, error) {
	vPart := v
	var err error
	for _, qe := range query {
		switch qe.kind {
		case kindKey:
			vPart, err = getKey(qe.key, vPart)
			if err != nil {
				return nil, err
			}
		default:
			return nil, UnknownError{
				Message: "invalid query",
			}
		}
	}
	return vPart, nil
}

func getKey(key string, v interface{}) (interface{}, error) {
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, ExpectedMapError{}
	}
	value, ok := m[key]
	if !ok {
		return nil, KeyNotFoundError{
			Key: key,
		}
	}
	return value, nil
}
