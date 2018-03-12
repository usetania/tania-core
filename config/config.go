package config

var Config Configuration

type Configuration struct {
	DemoMode       *bool
	UploadPathArea *string
	UploadPathCrop *string
	SqlitePath     *string
	MysqlHost      *string
	MysqlPort      *string
	MysqlDbname    *string
	MysqlUsername  *string
	MysqlPassword  *string
}
