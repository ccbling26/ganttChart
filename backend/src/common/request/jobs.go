package request

type JobsSearcher struct {
	ProductLines   []string `form:"productLines" json:"productLines" binding:"required"`
	StartTimeStamp int64    `form:"startTimeStamp" json:"startTimeStamp" binding:"required"`
	EndTimeStamp   int64    `form:"endTimeStamp" json:"endTimeStamp" binding:"required"`
	FaultTimeStamp int64    `form:"faultTimeStamp" json:"faultTimeStamp" binding:"required"`
}

// GetMessages 自定义错误信息
func (jobsSearcher JobsSearcher) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"ProductLines.required":   "产品线不能为空",
		"StartTimeStamp.required": "开始时间不能为空",
		"EndTimeStamp.required":   "结束时间不能为空",
		"FaultTimeStamp.required": "故障时间不能为空",
	}
}
