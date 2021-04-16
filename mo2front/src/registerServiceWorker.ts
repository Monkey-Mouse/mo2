/* eslint-disable no-console */

import axios from 'axios';
import { register } from 'register-service-worker'
import { Prompt, ShowRefresh } from './utils'

let offlinePrompt = false;

axios.interceptors.request.use((c) => {
  if (!offlinePrompt && !navigator.onLine) {
    Prompt('Oops, seems you are offline. We are now trying to serve you cached contents.', 50000);
    offlinePrompt = true;
  }
  return c;
})

if (process.env.NODE_ENV === 'production') {
  register(`https://www.motwo.cn/service-worker.js`, {
    ready() {
      console.log(
        'App is being served from cache by a service worker.\n' +
        'For more details, visit https://goo.gl/AFskqB'
      )
    },
    registered() {
      console.log('Service worker has been registered.')
    },
    cached() {
      console.log('Content has been cached for offline use.')
    },
    updatefound() {
      console.log('New content is downloading.')
    },
    updated() {
      console.log('New content is available; please refresh.')
      setTimeout(() => {
        ShowRefresh()
      }, 1000);
    },
    offline() {
      console.log('No internet connection found. App is running in offline mode.')
    },
    error(error) {
      console.error('Error during service worker registration:', error)
    }
  })
}
