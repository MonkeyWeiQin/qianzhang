import { GetAttachmentLabel, GetAttachmentList } from '@/api/admin'
export const commonConfig = {
  setAttachment :function(){
    GetAttachmentLabel().then((response) => {
      commonConfig.attachment = response.data
    })
  },
  setAttachmentList :function(){
    GetAttachmentList().then((response) => {
      commonConfig.attachmentList = response.data
    })
  },
  attachmentList: [],
  attachment: [],
}
