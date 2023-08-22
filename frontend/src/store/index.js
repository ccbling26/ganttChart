import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex)

const store = new Vuex.Store({
    state: {
        loggedIn: false,
    },
    mutations: {
        setLoggedIn(state, value) {
            state.loggedIn = value;
        }
    },
    getters: {
        getLoggedIn(state) {
            return state.loggedIn;
        }
    }
})

export default store
