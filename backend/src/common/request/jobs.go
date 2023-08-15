package request

type JobsSearcher struct {
	ProductLines   []string `form:"productLines" json:"productLines" binding:"required"`
	StartTimestamp int64    `form:"startTimestamp" json:"startTimestamp" binding:"required"`
	EndTimestamp   int64    `form:"endTimestamp" json:"endTimestamp" binding:"required"`
	FaultTimestamp int64    `form:"faultTimestamp" json:"faultTimestamp" binding:"required"`
}

// GetMessages 自定义错误信息
func (jobsSearcher JobsSearcher) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"ProductLines.required":   "产品线不能为空",
		"StartTimestamp.required": "开始时间不能为空",
		"EndTimestamp.required":   "结束时间不能为空",
		"FaultTimestamp.required": "故障时间不能为空",
	}
}
