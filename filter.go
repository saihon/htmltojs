package htmltojs

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Filter interface {
	Set(...atom.Atom)
	Add(...atom.Atom)
	Del(atom.Atom)
	len() int
	match(*html.Node) bool
}

type Data []atom.Atom

func (d *Data) Set(atoms ...atom.Atom) {
	*d = atoms
}

func (d *Data) Add(atoms ...atom.Atom) {
	*d = append(*d, atoms...)
}

func (d *Data) Del(atom atom.Atom) {
	for i, v := range *d {
		if atom == v {
			*d = append((*d)[:i], (*d)[i+1:]...)
			break
		}
	}
}

func (d Data) len() int {
	return len(d)
}

func (d *Data) match(node *html.Node) bool {
	for _, atom := range *d {
		if node.DataAtom == atom {
			return true
		}
	}
	return false
}

// Excludes
type Excludes struct {
	*Data
}

// Includes
type Includes struct {
	*Data
}

// Ignores
type Ignores struct {
	*Data
}
