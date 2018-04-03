package config

var Config Configuration

const (
	DB_INMEMORY = "inmemory"
	DB_SQLITE   = "sqlite"
	DB_MYSQL    = "mysql"
)

type Configuration struct {
	DemoMode               *bool
	UploadPathArea         *string
	UploadPathCrop         *string
	TaniaPersistenceEngine *string
	SqlitePath             *string
	MysqlHost              *string
	MysqlPort              *string
	MysqlDbname            *string
	MysqlUsername          *string
	MysqlPassword          *string
	RedirectURI            *string
	ClientID               *string
}
