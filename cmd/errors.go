/*
Copyright © 2023 orangekame3 miya.org.0309@gmai.com
*/
package cmd

import "errors"

var ErrNotFound = errors.New("record not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
