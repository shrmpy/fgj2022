
package acorn

import (

	//"log"
)

type Acorn struct {
	info func(string, ...any)
}
func NewAcorn(fn func(string, ...any)) *Acorn {
	ac := Acorn{ info: fn }
	return &ac
}
func (a *Acorn) Update() error {


	return nil
}


