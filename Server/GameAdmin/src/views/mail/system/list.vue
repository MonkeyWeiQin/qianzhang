<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <el-form-item label="邮件标题">
        <el-input placeholder="邮件标题" v-model="listQuery.uid"></el-input>
      </el-form-item>

      <el-form-item label="发送时间">
        <el-date-picker type="datetimerange"
                        v-model="listQuery.register_time"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        range-separator="至"
                        style="width: 100%;"
                        value-format="timestamp"
        ></el-date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="getSystemMailList()">查询</el-button>
        <el-button type="primary" @click="dialogFormCreate = true;sendAttachment=false;">发送普通邮件</el-button>
        <el-button type="primary" @click="dialogFormCreate = true;sendAttachment=true;">发送附件邮件</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="120" prop="mid" fixed="left" label="ID"/>
      <el-table-column width="300" prop="title" label="邮件标题"/>
      <el-table-column width="489" prop="content" label="邮件内容">
        <el-row slot-scope="scope">
          <div v-html="scope.row.content"></div>
        </el-row>
      </el-table-column>
      <el-table-column width="100" prop="uid" label="收件人ID">
        <template slot-scope="scope">
          <p v-if="scope.row.isAllUser === true">
            <el-tag :type="'success'" disable-transitions>全部发送</el-tag>
          </p>
          <p v-else>
            <el-tag :type="'success'" disable-transitions>uid</el-tag>
          </p>
        </template>
      </el-table-column>
      <el-table-column width="160" prop="sendTime" label="发送时间">
        <el-row slot-scope="scope">
          {{ scope.row.sendTime|formatterTime(that) }}
        </el-row>
      </el-table-column>
      <el-table-column width="160" prop="invalidTime" label="失效时间">
        <el-row slot-scope="scope">
          {{ scope.row.invalidTime|formatterTime(that) }}
        </el-row>
      </el-table-column>
      <el-table-column width="160" prop="cronTime" label="定时发送">
        <el-row slot-scope="scope">
          {{ scope.row.cronTime|formatterTime(that) }}
        </el-row>
      </el-table-column>

      <el-table-column width="150" label="操作">
        <el-row slot-scope="scope">
          <el-button round size="mini" type="text" @click="dialogFormCreate=true;
            MailForm.sendTime = scope.row.sendTime* 1000;
            MailForm.invalidTime= scope.row.invalidTime* 1000;
            MailForm.cronTime= scope.row.cronTime* 1000;
            MailForm.uid= scope.row.uid;
            MailForm.content= scope.row.content;
            MailForm.isAllUser= scope.row.isAllUser;
            MailForm.title= scope.row.title;
            MailForm.attachment = scope.row.attachment;
            MailForm.mid = scope.row.mid;
            sendAttachment = true"
          >编辑
          </el-button>
        </el-row>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination background layout="prev, pager, next" @current-change="currentChange" :page-size="listQuery.limit"
                     :total="total"
      />
    </div>
    <el-dialog title="新增/编辑权限" :visible.sync="dialogFormCreate">
      <el-form :model="MailForm" ref="MailForm" :rules="rules">
        <el-form-item label="发送时间" :label-width="formLabelWidth" prop="send_time">
          <el-date-picker type="datetime" v-model="MailForm.sendTime" value-format="timestamp"></el-date-picker>
        </el-form-item>
        <el-form-item label="失效时间" :label-width="formLabelWidth" prop="invalid_time">
          <el-date-picker type="datetime" v-model="MailForm.invalidTime" value-format="timestamp"></el-date-picker>
        </el-form-item>
        <el-form-item label="定时发送" :label-width="formLabelWidth" prop="cron_time">
          <el-time-picker style="width: 100%;" v-model="MailForm.cronTime" value-format="timestamp"></el-time-picker>
        </el-form-item>
        <div v-if="MailForm.isAllUser===false">
          <el-form-item label="收件人" :label-width="formLabelWidth" prop="is_all_user">
            <el-input v-model="MailForm.uid" autocomplete="off"></el-input>
          </el-form-item>
        </div>
        <el-form-item label="是否是全服发送" :label-width="formLabelWidth" prop="status">
          <el-switch
            v-model="MailForm.isAllUser"
            active-color="#13ce66"
            inactive-color="#ff4949"
          >
          </el-switch>
        </el-form-item>
        <el-form-item label="标题" :label-width="formLabelWidth" prop="title">
          <el-input v-model="MailForm.title" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="内容" :label-width="formLabelWidth" prop="content">
          <!--          <tinymce v-model="MailForm.content" :height="200"/>-->
          <el-input type="textarea" v-model="MailForm.content" autocomplete="off"></el-input>
        </el-form-item>
        <el-form v-if="sendAttachment===true" :inline="true" class="demo-dynamic form-inline">
          <div v-for="(domain, index) in MailForm.attachment">
            <el-form-item
              :label-width="formLabelWidth"
              :label="'附件' + (index+1)"
              :key="index"
              :prop="'attachment[' + index + '].itemType'"
            >
              <el-select v-model="domain.itemType" @change="ItemTypeChange(index)" placeholder="请选择附件类型">
                <el-option
                  v-for="item in itemTypeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item prop="itemId">
              <el-select v-model="domain.itemId" placeholder="请选择物品">
                <el-option
                  v-for="item in attachmentList[domain.itemType]"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                  <span style="float: left">{{ item.label }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ item.value }}</span>
                </el-option>
              </el-select>
            </el-form-item>
            <el-input-number style="width: 150px;" :min="1" prop="count" v-model="domain.count" label="描述文字"
            ></el-input-number>
            <el-button style="margin-left: 10px;" type="primary" @click="addDomain">新增附件</el-button>
            <el-button type="danger" @click.prevent="removeDomain(domain)">删除</el-button>
          </div>
        </el-form>
      </el-form>
      <div slot="default" class="dialog-footer" :label-width="formLabelWidth">
        <el-button @click="dialogFormCreate = false">取 消</el-button>
        <span v-if="MailForm.mid!==''">
          <el-button type="primary" @click="UpdateMail">确 定</el-button>
        </span>
        <span v-else>
          <el-button type="primary" @click="CreateMail">确 定</el-button>
        </span>
      </div>
    </el-dialog>
  </div>
</template>
<script>
// import Tinymce from '@/components/Tinymce'

import { getSystemMail, createSystemMail, updateSystemMail } from '@/api/system-mail'
import { commonConfig } from '@/common'

export default {
  name: 'SystemMailList',
  // components: {Tinymce},
  data() {
    return {
      rules: {},
      itemTypeOptions: commonConfig.attachment,
      attachmentList: commonConfig.attachmentList,
      dialogFormCreate: false,
      sendAttachment: false,
      formLabelWidth: '120px',
      list: [],
      total: 0,
      that: this,
      listQuery: {
        page: 1,
        limit: 5
      },
      MailForm: {
        mid: '',
        sendTime: '',
        invalidTime: '',
        cronTime: '',
        uid: '',
        isAllUser: false,
        title: '',
        content: '',
        attachment: [{
          itemType: null,
          itemId: null,
          count: null
        }]
      }
    }
  },
  created() {
    this.getSystemMailList()
  },

  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  methods: {
    getSystemMailList() {
      getSystemMail(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    addDomain() {
      this.MailForm.attachment.push({
        itemId: '',
        count: 0,
        itemType: ''
      })
    },
    removeDomain(item) {
      const index = this.MailForm.attachment.indexOf(item)
      if (index !== -1) {
        this.MailForm.attachment.splice(index, 1)
      }
    },
    currentChange(val) {
      this.listQuery.page = val
      this.getSystemMailList()
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },
    resetFields() {
      this.MailForm.attachment = [{
        itemType: null,
        itemId: null,
        count: null
      }]
      // Tinymce.methods.setContent("")
      this.$refs.MailForm.resetFields()
    },
    CreateMail() {
      if (this.MailForm.sendTime > 0) {
        this.MailForm.sendTime /= 1000
      }
      if (this.MailForm.cronTime > 0) {
        this.MailForm.cronTime /= 1000
      }
      if (this.MailForm.invalidTime > 0) {
        this.MailForm.invalidTime /= 1000
      }
      if (this.MailForm.uid.length === 0) {
        this.MailForm.uid = []
      } else {
        this.MailForm.uid = this.MailForm.uid.trim().split(',')
        this.MailForm.uid = this.MailForm.uid.map(item => {
          return +item
        })
      }
      createSystemMail(this.MailForm).then(() => {
        this.dialogFormCreate = false
        this.getSystemMailList()
        this.resetFields()
      })
    },
    UpdateMail() {
      if (this.MailForm.sendTime > 0) {
        this.MailForm.sendTime /= 1000
      }
      if (this.MailForm.cronTime > 0) {
        this.MailForm.cronTime /= 1000
      }
      if (this.MailForm.invalidTime > 0) {
        this.MailForm.invalidTime /= 1000
      }
      if (this.MailForm.uid.length === 0) {
        this.MailForm.uid = []
      } else {
        this.MailForm.uid = this.MailForm.uid.trim().split(',')
        this.MailForm.uid = this.MailForm.uid.map(item => {
          return +item
        })
      }
      updateSystemMail(this.MailForm).then(() => {
        this.dialogFormCreate = false
        this.getSystemMailList()
        this.resetFields()
      })
    },
    ItemTypeChange(index) {
      this.MailForm.attachment[index].itemId = null
    }
  }
}
</script>

<style scoped>
.tab-container {
  margin: 30px;
}
</style>
