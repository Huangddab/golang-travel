package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper // 存储Viper 实例指针
}

func NewSetting() (*Setting, error) {
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

	return &Setting{vp: vip}, nil
}
