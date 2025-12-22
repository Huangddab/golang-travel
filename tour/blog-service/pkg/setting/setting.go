package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper // 存储Viper 实例指针
}

// 热更新监听和变更处理

func NewSetting(configs ...string) (*Setting, error) {
	vip := viper.New() // 创建一个新的Viper实例
	// 查找名为 `config` 的配置文件（会查找 config.yaml）
	vip.SetConfigName("config")
	// 在 `configs` 目录下查找配置文件
	vip.AddConfigPath("configs")
	vip.SetConfigType("yaml") // 配置文件类型
	err := vip.ReadInConfig() // 读取配置文件
	if err != nil {
		return nil, err
	}
	s := &Setting{vp: vip}
	// 起一个协程 监听文件配置
	s.WatchSettingChange()

	return s, nil
}

func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
