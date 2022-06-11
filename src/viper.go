/**
 * @Author: YMBoom
 * @Description:
 * @File:  viper
 * @Version: 1.0.0
 * @Date: 2022/05/31 16:36
 */
package src

import (
	"embed"
	_ "embed"
	"github.com/spf13/viper"
	"strings"
)

//const confPath = "src/conf.yaml"

//go:embed  conf.yaml
var f embed.FS

func ReadConf() {
	data, err := f.ReadFile("conf.yaml")
	if err != nil {
		writeErr(err.Error())
		return
	}
	viper.SetConfigType("yaml")
	viper.ReadConfig(strings.NewReader(string(data)))
}

func GetMapStringMap(key string) map[string]string {
	return viper.GetStringMapString(key)
}
