package merge

import (
	"fmt"
	"github.com/reederc42/yt/errors"
	"strings"
)

func OrthogonalMerge(l, r interface{}) (interface{}, error) {
	return orthogonalMergeHelper(l, r, nil)
}

func orthogonalMergeHelper(l, r interface{},
	parents []string) (interface{}, error) {
	lMap, lOk := l.(map[interface{}]interface{})
	rMap, rOk := r.(map[interface{}]interface{})
	if !(lOk && rOk) {
		return nil, errors.KeyAlreadyDefined{
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
