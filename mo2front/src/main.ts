import './register.ts'
import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import VueCookies from 'vue-cookies'
import sanitizeHtml from 'sanitize-html'
import '@mdi/font/css/materialdesignicons.css'
const list = sanitizeHtml.defaults.allowedTags.concat(sanitizeHtml.defaults.selfClosing);
const attrs = sanitizeHtml.defaults.allowedAttributes;
attrs['*'] = ['href', 'align', 'alt', 'center', 'bgcolor', 'data-*', 'class', 'style']
const sanitize = (dirty) => {
  return sanitizeHtml(dirty, {
    allowedTags: list,
    allowedAttributes: attrs
  });
}
Vue.prototype.$sanitize = sanitize;
Vue.prototype.$d = Date.parse;
Vue.use(VueCookies);

// set default config
Vue.$cookies.config('7d')
Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
