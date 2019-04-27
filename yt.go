package yt

import (
	"gopkg.in/yaml.v2"
	"io"
	"regexp"
	"strings"

	ghodss "github.com/ghodss/yaml"
)

const inheritsRERaw = "#&inherits (.*)"
var inheritsRE *regexp.Regexp

func init() {
	inheritsRE = regexp.MustCompile(inheritsRERaw)
}

//only one document shall be in file
func Compile(input []byte) (interface{}, error) {
	var (
		err error
		v interface{}
	)
	parents := getParents(input)
	if len(parents) == 0 {
		err = yaml.Unmarshal(input, &v)
	} else {
		v = map[interface{}]interface{}{}
		for _, parent := range parents {
			p, queryErr := Query(nil, parent)
			if queryErr != nil {
				return nil, err
			}
			v, err = OrthogonalMerge(v, p)
			if err != nil {
				return nil, err
			}
		}
		var i interface{}
		err = yaml.Unmarshal(input, &i)
		v = Merge(v, i)
	}
	return v, err
}

func WriteYAML(v interface{}, output io.Writer) error {
	err := yaml.NewEncoder(output).Encode(v)
	return err
}

func WriteJSON(v interface{}, output io.Writer) error {
	rawYAML, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	rawJSON, err := ghodss.YAMLToJSON(rawYAML)
	if err != nil {
		return err
	}
	_, err = output.Write(rawJSON)
	return err
}

func getParents(doc []byte) []string {
	var parents []string
	for _, s := range strings.Split(string(doc), "\n") {
		m := inheritsRE.FindStringSubmatch(s)
		if len(m) >= 2 {
			parents = append(parents, m[1])
		}
	}
	return parents
}
