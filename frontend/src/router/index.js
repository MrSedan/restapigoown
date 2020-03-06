import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
      path: '/',
      name: "Home",
      component: () => import('../views/Home.vue')
    },
    {
      path: '/profile/:id(\\d+)',
      name: 'Profile',
      component: () => import('../views/Profile.vue')
    },
    {
      path: '/login',
      name: 'Log in',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/signup',
      name: "Sign Up",
      component: () => import('../views/Signup.vue')
    },
    {
      path: '/*',
      name: '404 not found',
      component: () => import('../views/404.vue')
    }
  // {
  //   path: '/about',
  //   name: 'About',
  //   route level code-splitting
  //   this generates a separate chunk (about.[hash].js) for this route
  //   which is lazy-loaded when the route is visited.
  //   component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  // }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
