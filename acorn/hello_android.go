
package acorn

import (

	//"log"
)

type Parcel struct {
	info func(string, ...any)
}
func NewParcel(fn func(string, ...any)) *Parcel {
	ac := Parcel{ info: fn }
	return &ac
}
func (p *Parcel) Update() error {


	return nil
}


