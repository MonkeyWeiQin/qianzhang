<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <el-form-item label="用户ID">
        <el-input placeholder="用户ID" v-model="listQuery.uid"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="listQuery.page=1;getMailList()">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="100" prop="uid" fixed="left" label="用户ID"/>
      <el-table-column width="200" prop="title" label="邮件标题"/>
      <el-table-column width="959" prop="content" label="邮件内容"/>
      <el-table-column width="80" prop="state" label="邮件状态">
        <template slot-scope="scope">
          <p v-if="scope.row.state === 0">
            <el-tag :type="'success'" disable-transitions>未读</el-tag>
          </p>
          <p v-if="scope.row.state === 1">
            <el-tag :type="'normal'" disable-transitions>已读</el-tag>
          </p>
          <p v-if="scope.row.state === 2">
            <el-tag :type="'danger'" disable-transitions>已删除</el-tag>
          </p>
        </template>
      </el-table-column>

      <el-table-column width="150" label="操作">
        <template slot-scope="scope">
          <el-popover
            placement="left"
            width="340"
            :ref="`popover-${scope.$index}`"
          >
            <el-table :data="attachmentData">
              <el-table-column width="100" property="itemType" label="类型">
                <template slot-scope="attachmentScope">
                  <el-tag v-if="attachmentScope.row.itemType===1">金币</el-tag>
                  <el-tag v-if="attachmentScope.row.itemType===2">钻石</el-tag>
                </template>
              </el-table-column>
              <el-table-column width="100" property="itemId" label="名称">
                <template slot-scope="attachmentScope">
                  <p v-if="attachmentScope.row.itemType===1">金币</p>
                  <p v-if="attachmentScope.row.itemType===2">钻石</p>
                </template>
              </el-table-column>
              <el-table-column width="100" property="count" label="数量"></el-table-column>
            </el-table>
            <el-button v-if="scope.row.attachment.length>0 " type="text" slot="reference" @click="attachmentData=scope.row.attachment">查看附件</el-button>
          </el-popover>
          <el-button v-if="scope.row.state!==2" type="text" @click="delFrom.mid = scope.row.mid;delFrom.uid = scope.row.uid;delFromDialog=true">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination background layout="prev, pager, next" @current-change="currentChange" :current-page="listQuery.page" :page-size="listQuery.limit"
                     :total="total"
      />
    </div>
    <el-dialog
      title="提示"
      :visible.sync="delFromDialog"
      width="30%"
      center
    >
      <span>确定删除</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="delFromDialog = false">取 消</el-button>
        <el-button type="primary" @click="delFromDialog = false;delMail()">确 定</el-button>
      </span>
    </el-dialog>

  </div>
</template>
<script>

import { DelMail, GetMail } from '@/api/mail'

export default {
  name: 'MailList',
  data() {
    return {
      attachmentData: [],
      formLabelWidth: '120px',
      delFromDialog: false,
      list: [],
      total: 0,
      that: this,
      listQuery: {
        page: 1,
        limit: 10,
        uid: ''
      },
      delFrom: {
        mid: '',
        uid: ''
      }
    }
  },
  created() {
    this.getMailList()
  },

  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },

  methods: {
    getMailList() {
      this.loading = true
      GetMail(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.getMailList()
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },
    delMail() {
      DelMail(this.delFrom).then(() => {
        this.delFromDialog = false
        this.getMailList()
      })
    }
  }
}
</script>

<style scoped>
.tab-container {
  margin: 30px;
}
</style>
