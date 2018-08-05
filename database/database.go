package database

type Database interface {
	Create(interface{}) error
	CleanDatabase()
}
