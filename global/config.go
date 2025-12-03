package global

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Config struct { //这个信息会帮助你生成V2ray/Clash/ShadowRocket的订阅链接,同时这个是互联网浏览器访问的地址
	AppPort                 string `desc:"app port" def:"80"`                                                                                //golang app 服务端口,可选,建议默认80或者443
	AllowUsers              string `desc:"allow users UUID" def:"903bcd04-79e7-429c-bf0c-0456c7de9cdc,903bcd04-79e7-429c-bf0c-0456c7de9cd1"` //单机模式下,允许的用户UUID
	IntervalSecond          string `desc:"interval second" def:"3600"`                                                                       //seconds 向主控服务器推送,流量使用情况的间隔时间
	RunAt                   string `desc:"run at" def:""`                                                                                    //optional run at
	EnableDataUsageMetering string `desc:"enable data usage metering" def:"true"`                                                            //是否开启用户流量统计,使用true 开启用户流量统计,使用false 关闭用户流量统计
	BufferSize              string `desc:"buffer size in bytes" def:"8192"`                                                                  //缓冲区大小,用于WebSocket和TCP/UDP读取
}

func (c Config) EnableUsageMetering() bool {
	return strings.ToLower(c.EnableDataUsageMetering) == "true"
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("0.0.0.0:%s", c.AppPort)
}
func (c Config) PushIntervalSecond() int {
	iv, err := strconv.ParseInt(c.IntervalSecond, 10, 32)
	if err != nil {
		log.Println("failed to parse interval second:", err)
		return 3600
	}
	return int(iv)
}

func (c Config) ListenPort() int {
	iv, err := strconv.ParseInt(c.AppPort, 10, 32)
	if err != nil {
		log.Println("failed to parse port:", err)
		return 80
	}
	return int(iv)
}

func (c Config) GetBufferSize() int {
	if c.BufferSize == "" {
		return 8192
	}
	iv, err := strconv.ParseInt(c.BufferSize, 10, 32)
	if err != nil {
		log.Println("failed to parse buffer size:", err)
		return 8192
	}
	return int(iv)
}

var cfg *Config

func (c Config) UserIDS() []string {
	parts := strings.Split(c.AllowUsers, ",")
	ids := make([]string, 0)
	for _, uid := range parts {
		uid = strings.TrimSpace(uid)
		if uid != "" {
			ids = append(ids, uid)
		}
	}
	return ids
}

func (c Config) PushInterval() time.Duration {
	if c.PushIntervalSecond() <= 0 {
		return time.Minute * 60
	}
	return time.Second * time.Duration(c.PushIntervalSecond())
}
