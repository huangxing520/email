package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.SugaredLogger

/*
setJSONEncoder 设置logger编码
*/
func setJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

/*
setLoggerWrite 设置logger写入文件
*/
func setLoggerWrite() zapcore.WriteSyncer {
	//create, _ := os.OpenFile("./test.log",os.O_CREATE|os.O_APPEND|os.O_RDWR,0744)
	create, err := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(err)
	}

	return zapcore.AddSync(create)
}
func init() {
	core := zapcore.NewCore(setJSONEncoder(), setLoggerWrite(), zap.InfoLevel)
	Logger = zap.New(core).Sugar()
}
