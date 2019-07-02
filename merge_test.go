package yt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrthogonalMerge(t *testing.T) {
	o1 := map[interface{}]interface{}{
		"key1": "value1",
	}
	o2 := map[interface{}]interface{}{
		"key2": "value2",
	}
	expected := map[interface{}]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	o3, err := OrthogonalMerge(o1, o2)
	assert.NoError(t, err)
	assert.Equal(t, expected, o3)
}

func TestOrthogonalMerge_KeyAlreadyDefinedError(t *testing.T) {
	o1 := map[interface{}]interface{}{
		"key1": "value1",
	}
	o2 := map[interface{}]interface{}{
		"key1": "value2",
	}
	expected := ErrKeyAlreadyDefined{
		Key: ".key1",
	}
	_, err := OrthogonalMerge(o1, o2)
	assert.Equal(t, expected, err)
}

func TestOrthogonalMerge_SubDocument(t *testing.T) {
	o1 := map[interface{}]interface{}{
		"o1": map[interface{}]interface{}{
			"key1": "value1",
		},
	}
	o2 := map[interface{}]interface{}{
		"o1": map[interface{}]interface{}{
			"key2": "value2",
		},
	}
	expected := map[interface{}]interface{}{
		"o1": map[interface{}]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}
	o3, err := OrthogonalMerge(o1, o2)
	assert.NoError(t, err)
	assert.Equal(t, expected, o3)
}

func TestOrthogonalMerge_SubDocumentKeyAlreadyDefined(t *testing.T) {
	o1 := map[interface{}]interface{}{
		"o1": map[interface{}]interface{}{
			"key1": map[interface{}]interface{}{
				"subkey1": "value1",
			},
		},
	}
	o2 := map[interface{}]interface{}{
		"o1": map[interface{}]interface{}{
			"key1": []interface{}{
				"value1",
				"value2",
			},
		},
	}
	expected := ErrKeyAlreadyDefined{
		Key: ".o1.key1",
	}
	_, err := OrthogonalMerge(o1, o2)
	assert.Equal(t, expected, err)
}

func TestOrthogonalMerge_InputsNotMaps(t *testing.T) {
	o1 := "value"
	o2 := []interface{}{
		"value1",
	}
	expected := ErrKeyAlreadyDefined{
		Key: ".",
	}
	_, err := OrthogonalMerge(o1, o2)
	assert.Equal(t, expected, err)
}

func TestMerge_SingleValue(t *testing.T) {
	o1 := map[interface{}]interface{}{
		"key1": "value1",
	}
	o2 := "value2"
	expected := "value2"
	o3 := Merge(o1, o2)
	assert.Equal(t, expected, o3)
}

func TestMerge_Objects(t *testing.T) {
	o1 := map[interface{}]interface{}{
		"key1": map[interface{}]interface{}{
			"key1": "value1",
		},
	}
	o2 := map[interface{}]interface{}{
		"key1": map[interface{}]interface{}{
			"key2": "value2",
		},
	}
	expected := map[interface{}]interface{}{
		"key1": map[interface{}]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}
	o3 := Merge(o1, o2)
	assert.Equal(t, expected, o3)
}
