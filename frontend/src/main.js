import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from 'axios'
import querystring from 'querystring'

Vue.config.productionTip = false
Vue.prototype.$http = axios
Vue.prototype.$qs = querystring

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
