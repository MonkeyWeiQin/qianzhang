import request from '@/utils/request'

export function userList(query) {
  return request({
    url: '/user/list',
    method: 'get',
    params: query
  })
}

export function modifyDiamond(query) {
  return request({
    url: '/user/modifyDiamond',
    method: 'post',
    data: query
  })
}
export function modifyGold(query) {
  return request({
    url: '/user/modifyGold',
    method: 'post',
    data: query
  })
}
export function modifyStrength(query) {
  return request({
    url: '/user/modifyStrength',
    method: 'post',
    data: query
  })
}
export function modifyStatus(query){
  return request({
    url: '/user/modifyStatus',
    method: 'post',
    data: query
  })
}
