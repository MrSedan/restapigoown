import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from 'axios'
import querystring from 'querystring'
import VueTitle from '@/components/VueTitle.js'

Vue.config.productionTip = false
Vue.prototype.$http = axios
Vue.prototype.$qs = querystring
Vue.component('vue-title', VueTitle)

String.prototype.capitalize = function(lower) {
  return (lower ? this.toLowerCase() : this).replace(/(?:^|\s)\S/g, function(a) { return a.toUpperCase(); });
};

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
