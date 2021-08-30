import request from '@/utils/request'


export function GetMail(query) {
  return request({
    url: '/user/mail/list',
    method: 'get',
    params: query
  })
}

export function DelMail(query) {
  return request({
    url: '/user/mail/del',
    method: 'post',
    data: query
  })
}
