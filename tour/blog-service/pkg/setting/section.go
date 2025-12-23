package setting

import "time"

// 用于声明配置属性的结构体
type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// app
type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	Upload                Upload
	DefaultContextTimeout time.Duration
}

// database
type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

// upload
type Upload struct {
	UploadSavePath       string
	UploadServerUrl      string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

// jwt
type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

// email
type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

// 针对重载应用配置 sections是配置存储
var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

// 重新读取配置
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
