<template>
  <el-container style="height: 100%; width: 100%">
    <el-header style="height: 5%; width:100%; text-align: center; background-color: #B3C0D1; color: #333;">
      甘特图
    </el-header>
    <el-container style="height: 95%; width: 100%">
      <el-aside style="height: 100%; width: 10%; background-color: rgb(238, 241, 246)">
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
      tooltip: {
        formatter: function (params) {
          const start = moment(params.value[1]).format("YYYY-MM-DD HH:mm:ss")
          const end = moment(params.value[2]).format("YYYY-MM-DD HH:mm:ss")
          return start + " - " + end;
        }
      },
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
          }
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
          type: "slider",
          show: true,
          xAxisIndex: 0,
          start: 20,
          end: 80,
          minSpan: 5,
        },
        {
          type: 'inside',
        },
        {
          type: "slider",
          show: true,
          yAxisIndex: 0,
          start: 20,
          end: 80,
          minSpan: 5,
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
          return time > new Date() || time < this.time_range[0] || time > this.time_range[1]
        },
        selectableRange: this.faultSelectableRange(this.time_range, this.fault_time),
      }
    },
  },
  methods: {
    isTheSameDay(start_time, end_time) {
      return start_time.getYear() === end_time.getYear() &&
          start_time.getMonth() === end_time.getMonth() &&
          start_time.getDate() === end_time.getDate()
    },
    faultSelectableRange(time_range, fault_time) {
      let start = "00:00:00"
      let end = "23:59:59"
      const start_time = time_range[0]
      const end_time = time_range[1]
      if (fault_time) {
        if (start_time && this.isTheSameDay(fault_time, start_time)) {
          start = moment(start_time).format("HH:mm:ss")
        }
        if (end_time && this.isTheSameDay(fault_time, end_time)) {
          end = moment(end_time).format("HH:mm:ss")
        } else if (!end_time && this.isTheSameDay(fault_time, new Date())) {
          end = moment(new Date()).format("HH:mm:ss")
        }
      }
      return start + " - " + end
    },
    renderChart(startTimestamp, endTimestamp) {
      this.chart.setOption({
        xAxis: {
          min: startTimestamp,
          max: endTimestamp,
          scale: true,
          axisLabel: {
            formatter: function (val) {
              val = Math.max(startTimestamp, val)
              return moment(val).format("YYYY-MM-DD HH:mm:ss");
            }
          },
          position: "top",
        },
        yAxis: {
          type: "category",
          data: this.ganttData.categories,
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
          },
        ],
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
      this.renderChart(startTimestamp, endTimestamp);
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
  console.log(rectShape)
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
