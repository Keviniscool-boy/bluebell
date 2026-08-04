import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios from './service/api'

Vue.prototype.$axios = axios;
Vue.config.productionTip = false


router.beforeEach((to, from, next) => {
  console.log(to);
  console.log(from);
  if (to.meta.requireAuth) {
    if (localStorage.getItem("loginResult")) {
      next();
    } else {
      if (to.path === '/login') {
        next();
      } else {
        next({
          path: '/login'
        })
      }
    }
  }
  else {
    next();
  }
  if (to.fullPath == "/login") {
    if (localStorage.getItem("loginResult")) {
      next({
        path: from.fullPath
      });
    } else {
      next();
    }
  }
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
