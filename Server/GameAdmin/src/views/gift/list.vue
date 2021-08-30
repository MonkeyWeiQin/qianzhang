<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <!--      <el-form-item label="用户名">-->
      <!--        <el-input placeholder="用户名" v-model="listQuery.username"></el-input>-->
      <!--      </el-form-item>-->
      <el-form-item>
        <!--        <el-button type="primary" @click="GetGiftList()">查询</el-button>-->
        <el-button type="primary" @click="dialogFormCreate = true;resetFields()">生成兑换码</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="50" prop="mid" label="ID"/>
      <el-table-column width="300" prop="name" label="名称"></el-table-column>

      <el-table-column width="100" prop="count" label="兑换次数"></el-table-column>
      <el-table-column width="150" prop="count" label="剩余兑换次数"></el-table-column>
      <el-table-column width="180" prop="status" label="有效期">
        <el-row slot-scope="scope">
          <span v-if="scope.row.effectiveTime>0">
            {{ scope.row.effectiveTime|formatterTime(that) }}
          </span>
          <span v-else>
            永久
          </span>
        </el-row>
      </el-table-column>
      <el-table-column width="300" prop="info" label="礼包说明"></el-table-column>

      <el-table-column width="150" label="操作">
        <el-row slot-scope="scope">
          <el-button type="text" size="mini" @click="dialogFormCreate = true;
           GiftForm.mid =scope.row.mid;
           GiftForm.name =scope.row.name;
           GiftForm.count =scope.row.count;
           GiftForm.info =scope.row.info;
           GiftForm.effectiveTime = scope.row.effectiveTime * 1000;
           GiftForm.attachment =scope.row.attachment;
          scope.row.effectiveTime>0?GiftForm.isForever = false:GiftForm.isForever=true;"
          >
            编辑
          </el-button>
          <el-button type="text" size="mini" @click="DownLoadGiftCode(scope.row.mid ,scope.row.name ); ">
            导出
          </el-button>
        </el-row>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination background layout="prev, pager, next" @current-change="currentChange" :page-size="listQuery.limit"
                     :total="total"
      />
    </div>
    <el-dialog title="新增/编辑" :visible.sync="dialogFormCreate">
      <el-form :model="GiftForm" ref="GiftForm" :rules="rules">
        <el-form-item label="名称" :label-width="formLabelWidth" required prop="name">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-input v-model="GiftForm.name" autocomplete="off"></el-input>
            </el-col>
          </el-row>

        </el-form-item>
        <el-form-item label="兑换个数" :label-width="formLabelWidth" required prop="count">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-input-number :min=1 v-model="GiftForm.count" autocomplete="off"
                               style="width: 100%;"
              ></el-input-number>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item label="是否永久" :label-width="formLabelWidth">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-switch v-model="GiftForm.isForever" active-color="#13ce66" inactive-color="#ff4949"></el-switch>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item label="有效期" :label-width="formLabelWidth" prop="effectiveTime" v-if="GiftForm.isForever===false">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-date-picker type="datetime" v-model="GiftForm.effectiveTime"
                              value-format="timestamp" style="width: 100%"
              ></el-date-picker>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item label="礼包说明" :label-width="formLabelWidth" prop="info">
          <el-row :gutter="24">
            <el-col :span="12">
              <el-input type="textarea" v-model="GiftForm.info" autocomplete="off"></el-input>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item>
          <el-row :gutter="24" v-for="(domain, index) in GiftForm.attachment" :key="index"
                  style="margin-left:0;margin-bottom: 10px;"
          >
            <el-col :span="8">
              <el-form-item
                :label-width="formLabelWidth"
                :label="'礼品' + (index+1)"
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
                    <span style="float: left">{{ item.label }}</span>
                    <span style="float: right; color: #8492a6; font-size: 13px">{{ item.value }}</span>
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="5">
              <el-form-item :prop="'attachment[' + index + '].itemId'">
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
            </el-col>
            <el-col :span="5">
              <el-input-number :prop="'attachment[' + index + '].count'" v-model="domain.count" label="描述文字"
                               style="width: 100%"
              ></el-input-number>
            </el-col>
            <el-col :span="2">
              <el-button type="primary" @click="addDomain">新增礼品</el-button>
            </el-col>
            <el-col :span="2" :offset=1>
              <el-button type="danger" @click.prevent="removeDomain(domain)">删除</el-button>
            </el-col>
          </el-row>
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormCreate = false">取 消</el-button>
        <span v-if="GiftForm.mid!==''">
          <el-button type="primary" @click="UpdateGiftCode">确 定</el-button>
        </span>
        <span v-else>
          <el-button type="primary" @click="CreateGiftCode">确 定</el-button>
        </span>
      </div>
    </el-dialog>

  </div>
</template>
<script>

import { CreateGiftCode, DownLoadGiftCode, GetGiftList } from '@/api/gift-code'
import { commonConfig } from '@/common'

export default {
  name: 'GiftList',
  data() {
    return {
      formLabelWidth: '120px',
      dialogFormCreate: false,
      that: this,
      itemTypeOptions: commonConfig.attachment,
      attachmentList: commonConfig.attachmentList,
      list: [],
      total: 0,
      options: [],
      rules: {},
      listQuery: {
        page: 1,
        limit: 5,
        username: ''
      },
      GiftForm: {
        mid: '',
        name: '',
        count: 1,
        info: '',
        isForever: false,
        effectiveTime: '',
        attachment: [{
          itemType: null,
          itemId: null,
          count: null
        }]
      },
      DeleteGiftForm: {
        mid: ''
      }
    }
  },
  created() {
    this.GetGiftList()
  },

  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  methods: {
    GetGiftList() {
      this.loading = true
      GetGiftList(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    resetFields() {
      this.$nextTick(() => {
        this.GiftForm.attachment = [{
          itemType: null,
          itemId: null,
          count: null
        }]
        this.$refs.GiftForm.resetFields()
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.GetGiftList()
    },

    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },
    addDomain() {
      this.GiftForm.attachment.push({
        itemId: null,
        count: null,
        itemType: null
      })
    },
    removeDomain(item) {
      const index = this.GiftForm.attachment.indexOf(item)
      if (index !== -1) {
        this.GiftForm.attachment.splice(index, 1)
      }
    },
    CreateGiftCode() {
      if (this.GiftForm.isForever) {
        this.GiftForm.effectiveTime = 0
      } else {
        this.GiftForm.effectiveTime /= 1000
      }
      CreateGiftCode(this.GiftForm).then(response => {
        this.GetGiftList()
        this.dialogFormCreate = false
        this.resetFields()
      })
    },
    UpdateGiftCode() {

    },
    DeleteGiftCode() {

    },
    DownLoadGiftCode(mid, name) {
      DownLoadGiftCode(mid, name)
    },
    ItemTypeChange(index) {
      this.GiftForm.attachment[index].itemId = null
    }
  }

}
</script>

<style scoped>
.tab-container {
  margin: 30px;
}
</style>
