import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list.vue'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: 'Dashboard', icon: 'dashboard', affix: true }
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: '/user',
    component: Layout,
    name: 'User',
    meta: {
      title: '用户管理',
      icon: 'user'
    },
    children: [
      {
        path: '/user/list',
        component: () => import('@/views/user/list'),
        name: 'user-management-list',
        meta: { title: '用户列表', icon: 'user' }
      }
    ]
  },
  {
    path: '/admin',
    component: Layout,
    name: 'Admin',
    meta: {
      title: '后台管理',
      icon: 'peoples'
    },
    children: [
      {
        path: '/admin/list',
        component: () => import('@/views/admin/list'),
        name: 'user-management-list',
        meta: { title: '后台账号列表', icon: 'list' }
      },
      {
        path: '/menu/list',
        component: () => import('@/views/menu/list'),
        name: 'menu-list',
        meta: { title: '后台菜单列表', icon: 'list' }
      }
    ]
  },
  {
    path: '/mail',
    component: Layout,
    name: 'Mail',
    meta: {
      title: '邮件管理',
      icon: 'email'
    },
    children: [
      {
        path: '/mail/list',
        component: () => import('@/views/mail/list'),
        name: 'mail-list',
        meta: { title: '用户邮件列表', icon: 'list' }
      },
      {
        path: '/mail/system/list',
        component: () => import('@/views/mail/system/list'),
        name: 'mail-system-list',
        meta: { title: '系统邮件', icon: 'list' }
      }
    ]
  },
  {
    path: '/gift',
    component: Layout,
    name: 'Gift',
    meta: {
      title: '兑换码管理',
      icon: 'el-icon-postcard'
    },
    children: [
      {
        path: '/gift/list',
        component: () => import('@/views/gift/list'),
        name: 'gift-list',
        meta: { title: '兑换码列表', icon: 'list' }
      },
      {
        path: '/gift/log/list',
        component: () => import('@/views/gift/log-list'),
        name: 'gift-log-list',
        meta: { title: '兑换记录', icon: 'list' }
      }
    ]
  },
  {
    path: '/商店管理',
    component: Layout,
    name: 'Shop',
    meta: {
      title: '商店管理',
      icon: 'el-icon-goods'
    },
    children: [
      {
        path: '/shop/purchase-goods/log',
        component: () => import('@/views/shop/purchase-goods/log'),
        name: 'shop-purchase-log',
        meta: { title: '商品购买记录', icon: 'list' }
      },
      {
        path: '/shop/chest-box/log',
        component: () => import('@/views/shop/chest-box/log'),
        name: 'shop-chest-log',
        meta: { title: '宝箱购买记录', icon: 'list' }
      }
    ]
  },
  {
    path: '/notice',
    component: Layout,
    name: 'Notice',
    meta: {
      title: '系统消息',
      icon: 'message'
    },
    children: [
      {
        path: '/notice/login-notice',
        component: () => import('@/views/notice/login-notice'),
        name: 'gift-list',
        meta: { title: '登录公告', icon: 'list' }
      },
      {
        path: '/notice/game-notice',
        component: () => import('@/views/notice/game-notice'),
        name: 'gift-log-list',
        meta: { title: '游戏内公告', icon: 'list' }
      },
      {
        path: '/notice/broadcast',
        component: () => import('@/views/notice/broadcast'),
        name: 'gift-log-list',
        meta: { title: '广播', icon: 'list' }
      }
    ]
  }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
