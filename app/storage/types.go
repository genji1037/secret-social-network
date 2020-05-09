package storage

// Connection represent mysql connection info.
type Connection struct {
	Host         string
	User         string
	Password     string
	Database     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}
