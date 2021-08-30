import request from '@/utils/request'

export function AdminList(query) {
  return request({
    url: '/admin/list',
    method: 'get',
    params: query
  })
}

export function AddAdmin(query) {
  return request({
    url: '/admin/create',
    method: 'post',
    data: query
  })
}

export function login(data) {
  return request({
    url: '/admin/login',
    method: 'post',
    data:data
  })
}

export function getInfo() {
  return request({
    url: '/admin/info',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/admin/logout',
    method: 'post'
  })
}

export function UpdatePassword(query) {
  return request({
    url: '/admin/updatePassword',
    method: 'post',
    data: query
  })
}

export function UpdateRole(query) {
  return request({
    url: '/admin/updateRole',
    method: 'post',
    data: query
  })
}

export function GetAttachmentLabel() {
  return request({
    url: '/system/getAttachmentLabel',
    method: 'get'
  })
}

export function GetAttachmentList() {
  return request({
    url: '/system/getAttachmentList',
    method: 'get'
  })
}
