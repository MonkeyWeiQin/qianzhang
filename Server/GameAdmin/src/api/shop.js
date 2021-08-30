import request from '@/utils/request'


export function GetGoodsList(query) {
  return request({
    url: '/system/getGoodsList',
    method: 'get',
    params: query
  })
}
export function GetPurchaseList(query) {
  return request({
    url: '/system/getGoodsPurchaseList',
    method: 'get',
    params: query
  })
}

export function GetChestList(query) {
  return request({
    url: '/system/getChestList',
    method: 'get',
    params: query
  })
}
export function GetChestPurchaseList(query) {
  return request({
    url: '/system/getChestPurchaseList',
    method: 'get',
    params: query
  })
}
export function GetChestContent(query) {
  return request({
    url: '/system/getChestContent',
    method: 'get',
    params: query
  })
}

