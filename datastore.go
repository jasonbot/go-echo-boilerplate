package boilerplateapi

// Datastore is an abstraction for talking to the database
type Datastore interface {
	LoadRecord(tableName string, record interface{}, primaryKeys ...string) error
	SaveRecord(tableName string, record interface{}, primaryKeys ...string) error
	DeleteRecord(tableName string, primaryKeys ...string) error
	Close() error
}
