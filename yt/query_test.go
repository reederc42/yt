package yt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLex_NoError(t *testing.T) {
	q := ".foo"
	expectedElements := []queryElement{
		{
			kind: kindKey,
			key: "foo",
		},
	}

	elements, err := lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
}

func TestLex_TwoElements(t *testing.T) {
	q := ".foo.bar"
	expectedElements := []queryElement{
		{
			kind: kindKey,
			key: "foo",
		},
		{
			kind: kindKey,
			key: "bar",
		},
	}

	elements, err := lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
}

func TestEmptyTokenError(t *testing.T) {
	q := ".."
	_, err := lex(q)
	assert.Equal(t, EmptyTokenError{}, err)
}

func TestExecQuery(t *testing.T) {
	m := map[interface{}]interface{}{
		"foo": "value",
	}
	qe := []queryElement{
		{
			kind: kindKey,
			key: "foo",
		},
	}
	v, err := execQuery(m, qe)
	assert.NoError(t, err)
	assert.Equal(t, "value", v)
}

func TestExecQuery_SubObject(t *testing.T) {
	m := map[interface{}]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": "value",
		},
	}
	qe := []queryElement{
		{
			kind: kindKey,
			key: "foo",
		},
		{
			kind: kindKey,
			key: "bar",
		},
	}

	v, err := execQuery(m, qe)
	assert.NoError(t, err)
	assert.Equal(t, "value", v)
}

func TestExecQuery_NotFound(t *testing.T) {
	m := map[interface{}]interface{}{
		"bar": "value",
	}
	qe := []queryElement{
		{
			kind: kindKey,
			key: "foo",
		},
	}

	_, err := execQuery(m, qe)
	assert.Equal(t, KeyNotFoundError{
		Key: "foo",
	}, err)
}
