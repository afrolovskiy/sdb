package sdb

// Storage stores variable values.
type Storage struct {
	vars map[string]string
}

// NewStorage returns new Storage instance
func NewStorage() *Storage {
	var st Storage
	st.vars = make(map[string]string)
	return &st
}

// Set sets the variable to the value.
func (st *Storage) Set(name string, value string) {
	st.vars[name] = value
}

// Get returns the value of the variable.
func (st *Storage) Get(name string) string {
	return st.vars[name]
}

// HasVariable returns true if the variable name exists.
func (st *Storage) HasVariable(name string) bool {
	_, ok := st.vars[name]
	return ok
}

// Unset removes the variable name from Storage instance.
func (st *Storage) Unset(name string) {
	if st.HasVariable(name) {
		delete(st.vars, name)
	}
}

// NumEqualTo returns the number of variables that are currently set to value.
func (st *Storage) NumEqualTo(value string) int {
	var count int
	for _, v := range st.vars {
		if v == value {
			count = count + 1
		}
	}
	return count
}

// Copy returns full copy of Storage instance.
func (st *Storage) Copy() *Storage {
	nst := NewStorage()
	for k, v := range st.vars {
		nst.vars[k] = v
	}
	return nst
}
