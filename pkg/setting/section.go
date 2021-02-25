package setting

import "time"

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	LogSavePath string
	LogFileName string
	LogFileExt  string
}

type DatabaseSettingS struct {
	DBType        string
	UserName      string
	Password      string
	Host          string
	DBName        string
	TablePrefix   string
	Charset       string
	ParseTime     bool
	MaxIdleConnes int
	MaxOpenConnes int
}

type EmailSettingS struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
	To       []string
}

var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		if err := s.ReadSection(k, v); err != nil {
			return err
		}
	}

	return nil
}
