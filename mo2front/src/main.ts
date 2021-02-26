import './register.ts'
import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import VueCookies from 'vue-cookies'
import sanitizeHtml from 'sanitize-html'
Vue.prototype.$sanitize = sanitizeHtml
Vue.use(VueCookies)

// set default config
Vue.$cookies.config('7d')
Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
