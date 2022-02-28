package response

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
}

func Result(success bool, code int, data interface{}, msg string) Response {
	// 开始时间
	return Response{
		success,
		code,
		data,
		msg,
	}
}

func Ok(data interface{}, msg string, code int) Response {
	return Result(true, code, data, msg)
}
func Fail(msg string, code int) Response {
	return Result(false, code, "", msg)
}

/*func OkWithMessage(message string, c *gin.Context) {
    Result(SUCCESS, map[string]interface{}{}, message, c)
}*/

//func OkWithData(data interface{}, c *gin.Context) {
//    Result(SUCCESS, data, "操作成功", c)
//}
//
//func OkWithDetailed(data interface{}, message string, c *gin.Context) {
//    Result(SUCCESS, data, message, c)
//}
//
//
//func FailWithMessage(message string, c *gin.Context) {
//    Result(ERROR, map[string]interface{}{}, message, c)
//}
//
//func FailWithDetailed(data interface{}, message string, c *gin.Context) {
//    Result(ERROR, data, message, c)
//}
