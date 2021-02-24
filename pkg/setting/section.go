package setting

import "time"

type ServerSetting struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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
