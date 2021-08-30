import request from "@/utils/request";
import axios from "axios";
import { getToken } from '@/utils/auth'

export function GetGiftList(query) {
  return request({
    url: '/system/getGiftList',
    method: 'get',
    params: query
  })
}

export function CreateGiftCode(params) {
  return request({
    url: '/system/createGift',
    method: 'post',
    data: params
  })
}

export function DownLoadGiftCode(params , name) {
  axios.get(process.env.VUE_APP_BASE_API+"/system/downLoadGiftList?mid="+params, {responseType: "blob",headers:{"X-token":getToken()}}).then(res => {
    let blobUrl = window.URL.createObjectURL(new Blob([res.data], {
      type: "application/vnd.ms-excel",
    }))
    const a = document.createElement("a");
    a.style.display = "none";
    a.download = name+".xlsx";
    a.href = blobUrl;
    a.click();
  });
}


export function GetGiftLogList(query) {
  return request({
    url: '/system/getGiftLogList',
    method: 'get',
    params: query
  })
}
