package yt

import (
	"gopkg.in/yaml.v2"
	"io"
)

//only returns first document from input
func Compile(input io.Reader) (interface{}, error) {
	var v interface{}
	err := yaml.NewDecoder(input).Decode(&v)
	return v, err
}

func Write(v interface{}, output io.Writer) error {
	err := yaml.NewEncoder(output).Encode(v)
	return err
}
