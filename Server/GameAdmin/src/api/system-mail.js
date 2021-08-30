import request from '@/utils/request'

export function getSystemMail(params) {
  return request({
    url: '/system/getSystemMail',
    method: 'get',
    params: params
  })
}

export function createSystemMail(params) {
  return request({
    url: '/system/createSystemMail',
    method: 'post',
    data: params
  })
}
export function updateSystemMail(params) {
  return request({
    url: '/system/updateSystemMail',
    method: 'post',
    data: params
  })
}


