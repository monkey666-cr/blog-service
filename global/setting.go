package global

import (
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
)

var (
	ServerSetting *setting.ServerSettingS
	AppSetting    *setting.AppSettingS
	EmailSetting  *setting.EmailSettingS
	Logger        *logger.Logger
)
