package yt

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestLex_NoError(t *testing.T) {
	q := ".foo"
	expectedElements := []queryElement{
		{
			kind: kindKey,
			key:  "foo",
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
			key:  "foo",
		},
		{
			kind: kindKey,
			key:  "bar",
		},
	}
	elements, err := lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
}

func TestLex_Index(t *testing.T) {
	q := ".1"
	expectedElements := []queryElement{
		{
			kind:  kindIndex,
			index: 1,
		},
	}
	elements, err := lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
}

func TestLex_File(t *testing.T) {
	q := "'file.yaml'"
	expectedElements := []queryElement{
		{
			kind: kindFile,
			file: "file.yaml",
		},
	}
	elements, err := lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
	q = "file.yaml.foo"
	expectedElements = []queryElement{
		{
			kind: kindFile,
			file: "file.yaml.foo",
		},
	}
	elements, err = lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
}

func TestLex_FileQuery(t *testing.T) {
	q := "'file.yaml'.foo"
	expectedElements := []queryElement{
		{
			kind: kindFile,
			file: "file.yaml",
		},
		{
			kind: kindKey,
			key:  "foo",
		},
	}
	elements, err := lex(q)
	assert.NoError(t, err)
	assert.Equal(t, expectedElements, elements)
}

func TestLex_EmptyTokenError(t *testing.T) {
	q := ".."
	_, err := lex(q)
	assert.Equal(t, ErrEmptyToken{}, err)
}

func TestExecQuery(t *testing.T) {
	m := map[interface{}]interface{}{
		"foo": "value",
	}
	qe := []queryElement{
		{
			kind: kindKey,
			key:  "foo",
		},
	}
	v, err := execQuery(m, qe, map[string]bool{})
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
			key:  "foo",
		},
		{
			kind: kindKey,
			key:  "bar",
		},
	}
	v, err := execQuery(m, qe, map[string]bool{})
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
			key:  "foo",
		},
	}
	_, err := execQuery(m, qe, map[string]bool{})
	assert.Equal(t, ErrKeyNotFound{
		Key: "foo",
	}, err)
}

func TestExecQuery_InsertKey(t *testing.T) {
	m := map[interface{}]interface{}{
		"bar": map[interface{}]interface{}{
			"baz": "value",
		},
	}
	insert := "value1"
	qe := []queryElement{
		{
			kind: kindKey,
			key:  "bar",
		},
		{
			kind: kindKey,
			key:  "baz",
		},
	}
	expected := map[interface{}]interface{}{
		"bar": map[interface{}]interface{}{
			"baz": "value1",
		},
	}
	v, err := insertQuery(m, insert, qe)
	assert.NoError(t, err)
	assert.Equal(t, expected, v)
}

func TestExecQuery_InsertIndex(t *testing.T) {
	m := []interface{}{
		"value0",
		"value1",
	}
	insert := "value2"
	qe := []queryElement{
		{
			kind:  kindIndex,
			index: 0,
		},
	}
	expected := []interface{}{
		"value2",
		"value1",
	}
	v, err := insertQuery(m, insert, qe)
	assert.NoError(t, err)
	assert.Equal(t, expected, v)
}

func TestGetKey(t *testing.T) {
	y := `key:
  value`
	var i interface{}
	err := yaml.Unmarshal([]byte(y), &i)
	assert.NoError(t, err)
	v, err := getKey("key", i)
	assert.NoError(t, err)
	assert.Equal(t, "value", v)
}

func TestGetIndex(t *testing.T) {
	y := `- value0
- value1`
	var i interface{}
	err := yaml.Unmarshal([]byte(y), &i)
	assert.NoError(t, err)
	v, err := getIndex(0, i)
	assert.NoError(t, err)
	assert.Equal(t, "value0", v)
}

func TestDetectCycle(t *testing.T) {
	_ = os.Chdir("testdata")
	cycle2, err := ioutil.ReadFile("cycle1.yaml")
	assert.NoError(t, err)
	expectedError := ErrCycleDetected{Source: "cycle2.yaml"}
	c := make(chan error, 1)
	go func(input []byte) {
		_, err := Compile(input, map[string]bool{})
		c <- err
	}(cycle2)
	select {
	case compileErr := <-c:
		assert.Equal(t, expectedError, compileErr)
	case <-time.After(2 * time.Second):
		assert.Fail(t, "cycle undetected after 2 seconds")
	}
}
