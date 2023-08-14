package apis

import (
	"backend/config"
	"backend/src/models"
	"backend/src/utils"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type IntermediateJobItem struct {
	priority int64
	values   [][]interface{}
}

type GanttData struct {
	Name  string        `json:"name"`
	Value []interface{} `json:"value"`
}

type JobsResult struct {
	Categories []string    `json:"categories"`
	Data       []GanttData `json:"data"`
}

type JobSlice [][]interface{}

func (j JobSlice) Len() int           { return len(j) }
func (j JobSlice) Less(x, y int) bool { return j[x][0].(int64) > j[y][0].(int64) }
func (j JobSlice) Swap(x, y int)      { j[x], j[y] = j[y], j[x] }

func calPriority(origin int64, startTimestamp int64, faultTimestamp int64, duration int64) int64 {
	if startTimestamp > faultTimestamp {
		return utils.Min(origin, duration+startTimestamp-faultTimestamp)
	} else {
		return utils.Min(origin, faultTimestamp-startTimestamp)
	}
}

func sortIntermediateJobData(data map[string]IntermediateJobItem) ([]string, []GanttData) {
	var intermediateData [][]interface{}
	for k, v := range data {
		item := []interface{}{v.priority, k}
		intermediateData = append(intermediateData, item)
	}
	sort.Sort(JobSlice(intermediateData))
	var categories []string
	var ganttData []GanttData
	for i, val := range intermediateData {
		name := val[1].(string)
		categories = append(categories, strings.Split(strings.Split(name, "_")[0], ".")[1])

		values := data[name].values
		for _, value := range values {
			value[0] = i
			ganttData = append(ganttData, GanttData{
				Name:  name,
				Value: value,
			})
		}
	}
	return categories, ganttData
}

func generateJobsData(productLines []string, startTimestamp int64, endTimestamp int64, faultTimestamp int64) JobsResult {
	duration := endTimestamp - startTimestamp

	jobsData := models.QueryJobs(
		config.Global.DB,
		productLines,
		time.Unix(startTimestamp/1000, 0),
		time.Unix(endTimestamp/1000, 0),
	)

	intermediateJobData := map[string]IntermediateJobItem{}
	for _, jobData := range jobsData {
		job := jobData.Job
		start := jobData.Start.Unix() * 1000
		end := jobData.UpdateTime.Unix() * 1000

		val, ok := intermediateJobData[job]
		if ok {
			val.priority = calPriority(val.priority, start, faultTimestamp, duration)
			val.values = append(
				val.values,
				[]interface{}{0, start, end},
			)
			intermediateJobData[job] = val
		} else {
			intermediateJobData[job] = IntermediateJobItem{
				priority: calPriority(math.MaxInt64, start, faultTimestamp, duration),
				values:   [][]interface{}{{0, start, end}},
			}
		}
	}
	categories, ganttData := sortIntermediateJobData(intermediateJobData)

	res := JobsResult{Categories: categories, Data: ganttData}
	return res
}

func QueryJobs(c *gin.Context) {
	res := map[string]interface{}{
		"code": 20001,
		"msg":  "查询成功",
		"data": map[string]interface{}{},
	}
	productLines, ok := c.GetQueryArray("productLines")
	startTimestamp, err1 := strconv.ParseInt(c.Query("startTimestamp"), 10, 64)
	endTimestamp, err2 := strconv.ParseInt(c.Query("endTimestamp"), 10, 64)
	faultTimestamp, err3 := strconv.ParseInt(c.Query("faultTimestamp"), 10, 64)
	if !ok {
		res["code"] = 20002
		res["msg"] = "缺少必填参数: 产品线"
	} else if err1 != nil {
		res["code"] = 20002
		res["msg"] = "开始时间解析失败"
	} else if err2 != nil {
		res["code"] = 20002
		res["msg"] = "结束时间解析失败"
	} else if err3 != nil {
		res["code"] = 20002
		res["msg"] = "故障时间解析失败"
	} else {
		data := generateJobsData(productLines, startTimestamp, endTimestamp, faultTimestamp)
		res["data"] = data
	}
	c.JSON(200, res)
}
