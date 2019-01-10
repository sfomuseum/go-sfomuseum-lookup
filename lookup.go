package lookup

type LookupWriter interface {
	Add(interface{}) error
}

type Lookup interface {
	Find(string) (interface{}, bool)
}
