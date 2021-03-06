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
      path: '/users',
      name: 'Users',
      component: () => import('../views/Users.vue')
    },
    {
      path: '/editprofile',
      name: 'Edit profile',
      component: () => import('@/views/EditProfile.vue')
    },
    {
      path: '/chat/:id(\\d+)',
      name: 'Chat',
      component: () => import('@/views/Chat.vue')
    },
    {
      path: '/*',
      name: '404 not found',
      component: () => import('../views/404.vue')
    }
]

const router = new VueRouter({
  mode: 'history',
  routes
})

export default router
