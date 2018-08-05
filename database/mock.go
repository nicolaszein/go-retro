package database

type Mock struct {
	Error error
}

func (m Mock) Create(interface{}) error {
	return m.Error
}
func (m Mock) CleanDatabase() {}
