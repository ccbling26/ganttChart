<template>
  <el-container style="height: 100%; width: 100%">
    <el-header style="height: 5%; width:100%; text-align: center; background-color: #409EFF; color: white;">
      <div style="height: 100%; width: 100%; align-items: center; justify-content: center; display: flex">
        甘特图
      </div>
    </el-header>
    <el-container style="height: 95%; width: 100%">
      <el-aside style="height: 100%; width: 10%; background-color: #F5F7FA">
        <el-menu default-active="1">
          <el-menu-item index="1" style="text-align: center">甘特图</el-menu-item>
        </el-menu>
      </el-aside>
      <el-main style="height: 100%; width: 90%">
        <el-row style="height: 10%; width: 100%">
          <el-col :span="5" :offset="1">
            <label for="startTime">产品线：</label>
            <el-select v-model="product_lines_selected" multiple placeholder="请选择">
              <el-option
                  v-for="item in product_lines_options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value">
              </el-option>
            </el-select>
          </el-col>
          <el-col :span="8">
            <label>时间范围：</label>
            <el-date-picker
                v-model="time_range"
                type="datetimerange"
                range-separator="至"
                :clearable="false"
                :picker-options="timeRangePickerOptions"
                :editable="false"
                format="yyyy-MM-dd HH:mm:ss">
            </el-date-picker>
          </el-col>
          <el-col :span="5">
            <label>故障时间：</label>
            <el-date-picker
                v-model="fault_time"
                type="datetime"
                :clearable="false"
                :picker-options="faultPickerOptions"
                :editable="false"
                format="yyyy-MM-dd HH:mm:ss">
            </el-date-picker>
          </el-col>
          <el-col :span="2">
            <el-button type="primary" round id="queryButton" @click="clickButton">查询</el-button>
          </el-col>
          <el-col :span="2">
            <el-button type="primary" round id="cleanButton" @click="resetButton">重置</el-button>
          </el-col>
        </el-row>
        <el-row id="ganttChart" style="height: 90%; width: 100%"></el-row>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import * as ECharts from "echarts";
import axios from "axios";
import moment from "moment";
import qs from "qs";

export default {
  name: "GanttChart",
  data() {
    const now = new Date()
    this.queryProductLines()
    return {
      product_lines_selected: [],
      product_lines_options: [],
      time_range: [new Date(now - 60 * 60 * 1000), now],
      fault_time: new Date(now - 30 * 60 * 1000),
      chart: null,
      ganttData: [],
    }
  },
  mounted() {
    this.chart = ECharts.init(document.getElementById("ganttChart"));
    const now = new Date()
    const initialOption = {
      tooltip: {},
      title: {
        text: "甘特图",
        left: "center",
      },
      xAxis: {
        min: now.getTime() - 24 * 60 * 60 * 1000,
        max: now.getTime(),
        scale: true,
        axisLabel: {
          formatter: function (val) {
            val = Math.max(now.getTime() - 24 * 60 * 60 * 1000, val)
            return moment(val).format("YYYY-MM-DD HH:mm:ss");
          },
          hideOverlap: true,
        },
        position: "top",
      },
      yAxis: {
        type: "category",
        data: [],
      },
      series: [],
      dataZoom: [
        {
          type: "inside",
          id: 'insideX',
          xAxisIndex: 0,
          start: 0,
          end: 100,
          zoomOnMouseWheel: false,
        },
        {
          type: "slider",
          id: 'sliderX',
          xAxisIndex: 0,
          start: 0,
          end: 100,
          zoomOnMouseWheel: false,
        },
        {
          type: "inside",
          id: 'insideY',
          yAxisIndex: 0,
          start: 50,
          end: 100,
          zoomOnMouseWheel: false,
          moveOnMouseWheel: true,
        },
        {
          type: 'slider',
          id: 'sliderY',
          yAxisIndex: 0,
          start: 50,
          end: 100,
          zoomOnMouseWheel: false,
        },
      ],
    }
    this.chart.setOption(initialOption);
  },
  computed: {
    timeRangePickerOptions() {
      return {
        disabledDate: (time) => {
          return time > new Date()
        },
      }
    },
    faultPickerOptions() {
      return {
        disabledDate: (time) => {
          return this.compareDate(time, new Date()) === 1 ||
              this.compareDate(this.time_range[0], time) === 1 ||
              this.compareDate(time, this.time_range[1]) === 1
        },
        selectableRange: this.faultSelectableRange(this.time_range, this.fault_time),
      }
    },
  },
  methods: {
    compareDate(time1, time2) {
      const year1 = time1.getFullYear()
      const month1 = time1.getMonth() + 1
      const day1 = time1.getDate()
      const year2 = time2.getFullYear()
      const month2 = time2.getMonth() + 1
      const day2 = time2.getDate()
      if (year1 > year2) {
        return 1
      } else if (year1 < year2) {
        return -1
      } else if (month1 > month2) {
        return 1
      } else if (month1 < month2) {
        return -1
      } else if (day1 > day2) {
        return 1
      } else if (day1 < day2) {
        return -1
      } else {
        return 0
      }
    },
    faultSelectableRange(time_range, fault_time) {
      let start = "00:00:00"
      let end = "23:59:59"
      const start_time = time_range[0]
      const end_time = time_range[1]
      if (fault_time) {
        if (start_time && this.compareDate(fault_time, start_time) === 0) {
          start = moment(start_time).format("HH:mm:ss")
        }
        if (end_time && this.compareDate(fault_time, end_time) === 0) {
          end = moment(end_time).format("HH:mm:ss")
        } else if (!end_time && this.compareDate(fault_time, new Date()) === 0) {
          end = moment(new Date()).format("HH:mm:ss")
        }
      }
      return start + " - " + end
    },
    renderChart(startTimestamp, endTimestamp, faultTimestamp) {
      this.chart.setOption({
        xAxis: {
          min: startTimestamp,
          max: endTimestamp,
          type: 'time',
          position: 'top',
          splitLine: {
            show: true,
            lineStyle: {
              color: ["gray"],
              width: 0.5,
              type: "dashed",
            }
          },
          axisLine: {
            show: false
          },
          axisTick: {
            lineStyle: {
              color: '#929ABA'
            }
          },
          axisLabel: {
            color: '#929ABA',
            inside: false,
            align: 'center',
            formatter: function (val) {
              val = Math.max(startTimestamp, val)
              return moment(val).format("YYYY-MM-DD HH:mm:ss");
            },
            hideOverlap: true,
          },
          scale: true,
        },
        yAxis: {
          type: "category",
          data: this.ganttData.categories,
          splitLine: {
            show: true,
            lineStyle: {
              color: ["gray"],
              width: 0.5,
              type: "solid",
            }
          },
        },
        series: [
          {
            type: "custom",
            renderItem: renderGanttItem,
            itemStyle: {
              opacity: 0.8
            },
            encode: {
              x: [1, 2],
              y: 0
            },
            data: this.ganttData.data,
            tooltip: {
              formatter: function (params) {
                const start = moment(params.value[1]).format("YYYY-MM-DD HH:mm:ss")
                const end = moment(params.value[2]).format("YYYY-MM-DD HH:mm:ss")
                return "产品线: " + params.name + "<br>开始时间: " + start + "<br>结束时间: " + end;
              }
            }
          },
          {
            type: 'line',
            name: "故障时间",
            animation: false,
            markLine: {
              data: [
                {
                  xAxis: faultTimestamp,
                  label: {
                    show: false
                  },
                  symbolKeepAspect: true,
                  lineStyle: {
                    color: "#ff0000",
                    width: 2,
                    type: "dashed",
                  },
                },
              ],
              symbol: ["none", "none"],
              tooltip: {
                formatter: function (params) {
                  return "故障时间<br>" + moment(params["data"]["xAxis"]).format("YYYY-MM-DD HH:mm:ss");
                }
              }
            },
            emphasis: {
              label: {
                show: false,
              }
            }
          },
        ],
        grid: {
          containLabel: true
        },
      });
    },
    async clickButton() {
      if (this.product_lines_selected.length === 0) {
        alert("至少选择一条产品线")
        return
      }
      if (!this.time_range || this.time_range.length !== 2) {
        alert("时间范围有误")
        return
      }
      if (!this.fault_time) {
        alert("故障时间不能为空")
        return
      }

      const startTimestamp = this.time_range[0].getTime()
      const endTimestamp = this.time_range[1].getTime()
      const faultTimestamp = this.fault_time.getTime()

      this.ganttData = await queryJobs(this.product_lines_selected, startTimestamp, endTimestamp, faultTimestamp);
      this.renderChart(startTimestamp, endTimestamp, faultTimestamp);
    },
    async resetButton() {
      const now = new Date().getTime()
      this.product_lines_options = []
      this.time_range = [new Date(now - 60 * 60 * 1000), new Date(now)]
      this.fault_time = new Date(now - 30 * 60 * 1000)
      this.chart = null
      this.ganttData = []
    },
    async queryProductLines() {
      try {
        const response = await axios.get("/api/product_lines");
        if (response.status !== 200) {
          alert("接口请求失败，状态码为" + response.status)
          return [];
        }
        const data = response.data;
        if (data.code !== 20001) {
          alert("接口请求失败，错误信息：" + data.msg);
          return [];
        }
        data.data.forEach(item => {
          this.product_lines_options.push({value: item.value, label: item.label});
        })
      } catch (error) {
        alert("接口请求异常")
        return [];
      }
    }
  }
};

function renderGanttItem(params, api) {
  const jobIndex = api.value(0);
  const start = api.coord([api.value(1), jobIndex]);
  const end = api.coord([api.value(2), jobIndex]);
  const height = api.size([0, 1])[1] * 0.6;
  const rectShape = ECharts.graphic.clipRectByRect(
      {
        x: start[0],
        y: start[1] - height / 2,
        width: end[0] - start[0],
        height: height
      },
      {
        x: params.coordSys.x,
        y: params.coordSys.y,
        width: params.coordSys.width,
        height: params.coordSys.height
      }
  );
  return (
      rectShape && {
        type: "rect",
        transition: ["shape"],
        shape: rectShape,
        style: api.style()
      }
  );
}

async function queryJobs(productLines, startTimestamp, endTimestamp, faultTimestamp) {
  try {
    const params = {
      startTimestamp: startTimestamp,
      endTimestamp: endTimestamp,
      faultTimestamp: faultTimestamp,
      productLines: productLines,
    }
    const response = await axios.get("/api/jobs", {
      params,
      paramsSerializer: params => qs.stringify(params, {arrayFormat: 'repeat'}),
    });
    if (response.status !== 200) {
      alert("接口请求失败，状态码为" + response.status)
      return [];
    }
    const data = response.data;
    if (data.code !== 20001) {
      alert("接口请求失败，错误信息：" + data.msg);
      return [];
    }
    return data.data;
  } catch (error) {
    alert("接口请求异常")
    return [];
  }
}

</script>

<style scoped>
</style>
