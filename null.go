package lookup

type NullLookup struct {
	Lookup
}

func NewNullLookup() (Lookup, error) {

	l := NullLookup{}
	return &l, nil
}

func (l *NullLookup) Find(path string) (interface{}, bool) {
	return nil, false
}
