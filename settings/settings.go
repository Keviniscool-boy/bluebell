package settings

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局配置实例
var Conf = new(Config)

// Config 所有配置的结构体
type Config struct {
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
}

// SnowflakeConfig 雪花算法配置
type SnowflakeConfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

func Init() error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")     // 配置文件名称 不需要带后缀
	viper.SetConfigType("yaml")       // 配置文件类型
	viper.AddConfigPath("./")         // 配置文件路径（项目根）
	viper.AddConfigPath("./settings") // 也在 settings 子目录查找配置文件
	err := viper.ReadInConfig()       // 读取配置文件
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 将配置反序列化到结构体
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return err
	}

	viper.WatchConfig() // 监控配置文件变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改")
		// 配置文件变更后重新加载
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})

	return nil
}
