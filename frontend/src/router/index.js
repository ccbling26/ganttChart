import Vue from "vue";
import VueRouter from "vue-router";
import Login from "../components/FeishuLogin.vue";
import GanttChart from "../components/GanttChart.vue";
import store from "../store";

Vue.use(VueRouter)

const routes = [{
    path: "/", name: "index", component: Login,
}, {
    path: "/login", name: "login", component: Login,
}, {
    path: "/gantt_chart", name: "gantt_chart", component: GanttChart, meta: {
        requiresAuth: true,
    }
}]

const router = new VueRouter({
    routes
})

router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (!store.state.loggedIn) {
            next({
                path: "/login", query: {redirect: to.fullPath}
            });
        } else {
            next();
        }
    } else {
        next();
    }
})

export default router
