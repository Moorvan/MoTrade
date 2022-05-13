package core

import (
	"MoTrade/global"
	mlog "MoTrade/mo-log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper(path string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		mlog.Log.Fatalln("Fail to read config:", err.Error())
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		mlog.Log.Println("Config file changed:", e.Name)
		if err := v.Unmarshal(&global.GB_CONFIG); err != nil {
			mlog.Log.Errorln("Fail to update config:", err.Error())
		}
	})

	if err := v.Unmarshal(&global.GB_CONFIG); err != nil {
		mlog.Log.Fatalln("Fail to get config: " + err.Error())
	}

	return v
}
