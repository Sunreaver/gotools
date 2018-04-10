package database

// Configer Configer
type Configer interface {
	GetDBName() string
	GetCollectionName() string
	GetSocketTimeoutSecond() int
}
