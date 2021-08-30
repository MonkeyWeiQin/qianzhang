import request from '@/utils/request'


export function GetMenuList(query) {
  return request({
    url: '/system/getMenuList',
    method: 'get',
    params: query
  })
}

export function CreateMenu(query) {
  return request({
    url: '/system/createMenuData',
    method: 'post',
    data: query
  })
}

export function UpdateMenu(query) {
  return request({
    url: '/system/updateMenuData',
    method: 'post',
    data: query
  })
}

export function GetLevelMenuList(query) {
  return request({
    url: '/system/getLevelMenuList',
    method: 'get',
    params: query
  })
}

export function DeleteMenu(query) {
  return request({
    url: '/system/deleteMenuData',
    method: 'post',
    data: query
  })
}
