import request from '@/utils/request'

export function getBroadcastList(query) {
  return request({
    url: '/system/getBroadcast',
    method: 'get',
    params: query
  })
}

export function createBroadcast(params) {
  return request({
    url: '/system/createBroadcast',
    method: 'post',
    data: params
  })
}

export function updateBroadcast(params) {
  return request({
    url: '/system/updateBroadcast',
    method: 'post',
    data: params
  })
}
