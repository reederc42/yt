package query

import (
	"github.com/reederc42/yt/errors"
	"strconv"
	"unicode"
)

const (
	kindKey int = iota
	kindIndex
)

type state int
const (
	invalid state = iota
	start
	key
	index
)

type queryElement struct {
	kind int
	key string
	index int
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
				return nil, errors.EmptyToken{}
			case key:
				qe = append(qe, queryElement{
					kind: kindKey,
					key: token,
				})
				s = start
				token = ""
			case index:
				index, err := strconv.ParseInt(token, 10, 0)
				if err != nil {
					return nil, err
				}
				qe = append(qe, queryElement{
					kind: kindIndex,
					index: int(index),
				})
				s = start
				token = ""
			}
		case unicode.IsNumber(r):
			switch s {
			case invalid:
				return nil, errors.UnexpectedChar{
					Char: r,
				}
			case start:
				token = string(r)
				s = index
			case key:
				token += string(r)
			case index:
				token += string(r)
			}
		default:
			switch s {
			case invalid:
				return nil, errors.UnexpectedChar{
					Char: r,
				}
			case start:
				token = string(r)
				s = key
			case key:
				token += string(r)
			case index:
				return nil, errors.UnexpectedChar{
					Char: r,
				}
			}
		}
	}

	switch s {
	case invalid:
		return nil, errors.Unknown{
			Message: "invalid query",
		}
	case key:
		qe = append(qe, queryElement{
			kind: kindKey,
			key: token,
		})
	case index:
		index, err := strconv.ParseInt(token, 10, 0)
		if err != nil {
			return nil, err
		}
		qe = append(qe, queryElement{
			kind: kindIndex,
			index: int(index),
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
		case kindIndex:
			vPart, err  = getIndex(qe.index, vPart)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.Unknown{
				Message: "invalid query",
			}
		}
	}
	return vPart, nil
}

func getKey(key string, v interface{}) (interface{}, error) {
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, errors.ExpectedMap{}
	}
	value, ok := m[key]
	if !ok {
		return nil, errors.KeyNotFound{
			Key: key,
		}
	}
	return value, nil
}

func getIndex(index int, v interface{}) (interface{}, error) {
	a, ok := v.([]interface{})
	if !ok {
		return nil, errors.ExpectedArray{}
	}
	if index >= len(a) || index < 0 {
		return nil, errors.OutOfBounds{}
	}
	value := a[index]
	return value, nil
}
