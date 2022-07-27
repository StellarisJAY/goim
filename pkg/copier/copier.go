package copier

import "github.com/jinzhu/copier"

func CopyStructFields(to interface{}, from interface{}) error {
	return copier.Copy(to, from)
}
