package yt

import (
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"

	"github.com/reederc42/yt/errors"
)

const (
	kindKey int = iota
	kindIndex
	kindFile
	kindValue
)

type stateVal int

const (
	invalid stateVal = iota
	startElement
	key
	index
	file
)

type queryElement struct {
	kind  int
	key   string
	index int
	file  string
}

func Query(v interface{}, query string) (interface{}, error) {
	elements, err := lex(query)
	if err != nil {
		return nil, err
	}
	value, err := execQuery(v, elements)
	return value, err
}

func Insert(src, value interface{}, query string) (interface{}, error) {
	elements, err := lex(query)
	if err != nil {
		return nil, err
	}
	v, err := insertQuery(src, value, elements)
	return v, err
}

func lex(query string) ([]queryElement, error) {
	state := invalid
	var token string
	qe := make([]queryElement, 0)
	for _, r := range query {
		switch {
		case r == '.':
			switch state {
			case invalid:
				state = startElement
			case startElement:
				return nil, errors.EmptyToken{}
			case key:
				qe = append(qe, queryElement{
					kind: kindKey,
					key:  token,
				})
				state = startElement
				token = ""
			case index:
				index, err := strconv.ParseInt(token, 10, 0)
				if err != nil {
					return nil, err
				}
				qe = append(qe, queryElement{
					kind:  kindIndex,
					index: int(index),
				})
				state = startElement
				token = ""
			case file:
				token += string(r)
			}
		case unicode.IsNumber(r):
			switch state {
			case invalid:
				token += string(r)
				state = file
			case startElement:
				token = string(r)
				state = index
			case key:
				token += string(r)
			case index:
				token += string(r)
			case file:
				token += string(r)
			}
		case r == '\'' || r == '"':
			switch state {
			case invalid:
				token += string(r)
				state = file
			case startElement:
				fallthrough
			case key:
				fallthrough
			case index:
				return nil, errors.UnexpectedRune{
					Rune: r,
				}
			case file:
				fileName := strings.TrimLeft(strings.TrimRight(token, `'"`), `'"`)
				qe = append(qe, queryElement{
					kind: kindFile,
					file: fileName,
				})
				state = invalid
				token = ""
			}
		default:
			switch state {
			case invalid:
				token += string(r)
				state = file
			case startElement:
				token = string(r)
				state = key
			case key:
				token += string(r)
			case index:
				return nil, errors.UnexpectedRune{
					Rune: r,
				}
			case file:
				token += string(r)
			}
		}
	}

	switch state {
	case invalid:
		if len(qe) == 0 {
			return nil, errors.InvalidQuery{}
		}
	case key:
		qe = append(qe, queryElement{
			kind: kindKey,
			key:  token,
		})
	case index:
		index, err := strconv.ParseInt(token, 10, 0)
		if err != nil {
			return nil, err
		}
		qe = append(qe, queryElement{
			kind:  kindIndex,
			index: int(index),
		})
	case file:
		fileName := strings.TrimLeft(strings.TrimRight(token, `'"`), `'"`)
		qe = append(qe, queryElement{
			kind: kindFile,
			file: fileName,
		})
	}
	return qe, nil
}

//Parsing is implied by query elements. If the first query element is a file,
// the value of v is ignored.
func execQuery(v interface{}, query []queryElement) (interface{}, error) {
	vPart := v
	var err error
	for i, qe := range query {
		switch qe.kind {
		case kindKey:
			vPart, err = getKey(qe.key, vPart)
			if err != nil {
				return nil, err
			}
		case kindIndex:
			vPart, err = getIndex(qe.index, vPart)
			if err != nil {
				return nil, err
			}
		case kindFile:
			if i > 0 {
				return nil, errors.InvalidQuery{}
			}
			f, fileError := ioutil.ReadFile(qe.file)
			if fileError != nil {
				return nil, fileError
			}
			vPart, err = Compile(f)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.InvalidQuery{}
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

//Recursive insert; does not support file kinds
func insertQuery(src, insert interface{}, qe []queryElement) (interface{}, error) {
	if len(qe) == 0 {
		return insert, nil
	}
	var err error
	q := qe[0]
	switch q.kind {
	case kindKey:
		m, ok := src.(map[interface{}]interface{})
		if !ok {
			return nil, errors.ExpectedMap{}
		}
		m[q.key], err = insertQuery(m[q.key], insert, qe[1:])
		return m, err
	case kindIndex:
		a, ok := src.([]interface{})
		if !ok {
			return nil, errors.ExpectedArray{}
		}
		if q.index < 0 || q.index >= len(a) {
			return nil, errors.OutOfBounds{}
		}
		a[q.index], err = insertQuery(a[q.index], insert, qe[1:])
		return a, err
	default:
		return nil, errors.Unknown{Message: "insert query: bad kind"}
	}
}
