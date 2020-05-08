package storage

// 连接信息
type Connection struct {
	Host         string
	User         string
	Password     string
	Database     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}
