import request from '@/utils/request'

export function getNoticeList(query) {
  return request({
    url: '/system/getNotice',
    method: 'get',
    params: query
  })
}

export function createNotice(params) {
  return request({
    url: '/system/addNotice',
    method: 'post',
    data: params
  })
}
export function updateNotice(params) {
  return request({
    url: '/system/updateNotice',
    method: 'post',
    data: params
  })
}
