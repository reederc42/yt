package yt

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParents(t *testing.T) {
	_ = os.Chdir("testdata")
	f, err := ioutil.ReadFile("parent.yaml")
	assert.NoError(t, err)
	parents := getParents(f)
	expected := []string{
		"'list.yaml'.items.1",
	}
	assert.Equal(t, expected, parents)
}

func TestCompile(t *testing.T) {
	_ = os.Chdir("testdata")
	f, err := ioutil.ReadFile("service.yaml")
	assert.NoError(t, err)
	expected := map[interface{}]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[interface{}]interface{}{
			"name": "my-service",
		},
		"spec": map[interface{}]interface{}{
			"selector": map[interface{}]interface{}{
				"app": "nginx",
			},
			"ports": []interface{}{
				map[interface{}]interface{}{
					"protocol": "TCP",
					"port":     80,
				},
			},
		},
	}

	m, err := Compile(f, map[string]bool{})
	assert.NoError(t, err)
	assert.Equal(t, expected, m)
}

func TestCompile_WithParent(t *testing.T) {
	_ = os.Chdir("testdata")
	f, err := ioutil.ReadFile("parent.yaml")
	expected := map[interface{}]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[interface{}]interface{}{
			"name": "my-secure-service",
		},
		"spec": map[interface{}]interface{}{
			"selector": map[interface{}]interface{}{
				"app": "nginx",
			},
			"ports": []interface{}{
				map[interface{}]interface{}{
					"protocol": "TCP",
					"port":     8080,
				},
			},
		},
	}
	m, err := Compile(f, map[string]bool{})
	assert.NoError(t, err)
	assert.Equal(t, expected, m)
}
