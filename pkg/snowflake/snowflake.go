package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node //定义全局变量
func Init(startTime string, machineID int64) (err error) {
	var st time.Time                                       //定义一个时间变量
	st, err = time.Parse("2006-01-02 15:04:05", startTime) //将字符串类型的时间转换为time.Time类型
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1e6    //将时间转换为毫秒
	node, err = sf.NewNode(machineID) //创建一个新的雪花节点
	return
}
func GenerateID() int64 {
	return node.Generate().Int64() //生成一个新的雪花ID
}
