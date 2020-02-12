package htmltojs

const (
	_A = 97
	_Z = 122
)

// Variable is generate variable name.
type Variable struct {
	// prefix of variable name. that for an avoid the reserved words.
	Prefix string
	name   []byte
}

func (v *Variable) toString() string {
	return v.Prefix + string(v.name)
}

// GenName is returns generated a new name.
func (v *Variable) genName() string {
	if len(v.name) == 0 {
		v.name = append(v.name, _A)
		return v.toString()
	}

	for i := len(v.name) - 1; i >= 0; i-- {
		if v.name[i] < _Z {
			v.name[i]++
			break
		}

		v.name[i] = _A

		if i == 0 {
			v.name = append(v.name, _A)
			break
		}
	}

	return v.toString()
}
