<template>
  <div></div>
</template>
<script>
export default {
  name: "FeishuLogin",
  mounted() {
    if (this.apiAuth()) {
      this.$store.state.loggedIn = true
      this.$router.push("/gantt_chart")
    }
  },
  methods: {
    apiAuth() {
      if (!window.h5sdk) {
        alert("只能在飞书上运行");
        return;
      }
      const url = encodeURIComponent(location.href.split("#")[0]);
      fetch(`/api/get_config_parameters?url=${url}`).then((response) => response.json()
          .then((res) => {
            res = res.data
            console.log(res)
            // 调用config接口进行鉴权
            window.h5sdk.config({
              appId: res.appid,
              timestamp: res.timestamp,
              nonceStr: res.nonceStr,
              signature: res.signature,
              jsApiList: [],
              onSuccess: (res) => {
                console.log(`config success: ${JSON.stringify(res)}`);
              },
              onFail: (err) => {
                throw `config failed: ${JSON.stringify(err)}`;
              }
            });
            window.h5sdk.ready(() => {
              window.tt.getUserInfo({
                success(res) {
                  console.log(res)
                },
                fail(err) {
                  throw `获取用户信息失败: ${JSON.stringify(err)}`
                }
              })
            })
          })
      ).catch(function (e) {
        console.log(e)
        return false
      });
      return true
    }
  }
}
</script>
<style>

</style>