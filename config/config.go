package config

var Config Configuration

type Configuration struct {
	DemoMode       *bool
	UploadPathArea *string
	UploadPathCrop *string
}
