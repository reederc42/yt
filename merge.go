package yt

import (
	"fmt"
	"strings"
)

func Merge(l, r interface{}) interface{} {
	lMap, lOk := l.(map[interface{}]interface{})
	rMap, rOk := r.(map[interface{}]interface{})
	if !(lOk && rOk) {
		return r
	}
	for k, v := range rMap {
		lMap[k] = Merge(lMap[k], v)
	}
	return lMap
}

func OrthogonalMerge(l, r interface{}) (interface{}, error) {
	return orthogonalMergeHelper(l, r, nil)
}

func orthogonalMergeHelper(l, r interface{},
	parents []string) (interface{}, error) {
	lMap, lOk := l.(map[interface{}]interface{})
	rMap, rOk := r.(map[interface{}]interface{})
	if !(lOk && rOk) {
		return nil, ErrKeyAlreadyDefined{
			Key: fmt.Sprintf(".%s", strings.Join(parents, ".")),
		}
	}

	var err error
	for rk, rv := range rMap {
		lv, ok := lMap[rk]
		if !ok {
			lMap[rk] = rv
		} else {
			lMap[rk], err = orthogonalMergeHelper(lv, rv,
				append(parents, fmt.Sprintf("%v", rk)))
			if err != nil {
				return nil, err
			}
		}
	}

	return lMap, nil
}
