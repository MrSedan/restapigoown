import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from 'axios'
import querystring from 'querystring'
import VueTitle from '@/components/VueTitle.vue'

Vue.config.productionTip = false
Vue.prototype.$http = axios
Vue.prototype.$qs = querystring
Vue.component('vue-title', VueTitle)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
