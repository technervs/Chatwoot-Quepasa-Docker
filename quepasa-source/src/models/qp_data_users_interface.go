package models

type QpDataUsersInterface interface {
	Count() (int, error)
	Find(string) (*QpUser, error)
	Exists(string) (bool, error)
	Check(string, string) (*QpUser, error)
	Create(string, string) (*QpUser, error)
}
