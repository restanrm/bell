import Vue from 'vue'
import Router from 'vue-router'
import List from '@/components/List'
import SoundUpload from '@/components/SoundUpload'
import destination from '@/components/Destination'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'List',
      component: List
    },
    {
      path: '/upload',
      name: 'SoundUpload',
      component: SoundUpload
    },
    {
      path: '/destination',
      name: 'destination',
      component: destination
    }
  ]
})
