<template>
  <body>
  <div id="app">
    <GanttChart/>
  </div>
  </body>
</template>

<script>
import GanttChart from './components/GanttChart.vue'

export default {
  name: 'App',
  components: {
    GanttChart
  }
}
</script>

<style>
html,
body,
#app {
  height: 100%;
  width: 100%;
  margin: 0 auto;
  padding: 0;
}
</style>
