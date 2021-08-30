<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <el-form-item label="使用人ID">
        <el-input placeholder="使用人ID" v-model="listQuery.uid"></el-input>
      </el-form-item>
      <el-form-item label="兑换码">
        <el-input placeholder="兑换码" v-model="listQuery.code"></el-input>
      </el-form-item>
      <el-form-item label="兑换码名称">
        <el-select v-model="listQuery.mid" placeholder="请选择">
          <el-option label="全部" value=""></el-option>
          <el-option
            v-for="item in GiftOptions"
            :key="item.value"
            :label="item.name"
            :value="item.mid">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="是否被使用">
        <el-select v-model="listQuery.used" placeholder="请选择">
          <el-option
            v-for="item in UsedOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="listQuery.page=1;GetGiftLogList();">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="100" prop="code" label="兑换码"/>
      <el-table-column width="180" prop="uid" label="是否被使用">
        <el-row slot-scope="scope">
          <span v-if="scope.row.uid>0">
            已使用
          </span>
          <span v-else>
            未使用
          </span>
        </el-row>
      </el-table-column>

      <el-table-column width="160" label="兑换码名称" prop="mid">
        <el-row slot-scope="scope">
          {{ scope.row.mid|findGiftName(that) }}
        </el-row>
      </el-table-column>
      <el-table-column width="200" label="兑换有效期" prop="mid">
        <el-row slot-scope="scope">
          {{ scope.row.mid|findGiftEffectiveTime(that) }}
        </el-row>
      </el-table-column>
      <el-table-column width="400" label="礼包说明" prop="mid">
        <el-row slot-scope="scope">
          {{ scope.row.mid|findGiftInfo(that) }}
        </el-row>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination background layout="prev, pager, next" @current-change="currentChange" :page-size="listQuery.limit"
                     :total="total"/>
    </div>
  </div>
</template>
<script>

import {GetGiftLogList, GetGiftList} from '@/api/gift-code'

export default {
  name: 'GiftLogList',
  data() {
    return {
      that: this,
      UsedOptions: [
        {
          value: 0,
          label: '全部'
        },
        {
          value: 1,
          label: '已使用'
        },
        {
          value: 2,
          label: '未使用'
        }
      ],
      GiftOptions: [],
      list: [],
      total: 0,
      listQuery: {
        page: 1,
        code: "",
        limit: 15,
        uid: "",
        used: "",
        mid: "",
      },
    }
  },
  created() {
    GetGiftList({limit: 999999}).then(response => {
      this.GiftOptions = response.data.list
    })
    this.GetGiftLogList()
  },
  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    },
    findGiftInfo: (value, that) => {
      for (let i = 0; i < that.GiftOptions.length; i++) {
        if (that.GiftOptions[i].mid === value) {
          return that.GiftOptions[i].info
        }
      }
    },
    findGiftName: (value, that) => {
      for (let i = 0; i < that.GiftOptions.length; i++) {
        if (that.GiftOptions[i].mid === value) {
          return that.GiftOptions[i].name
        }
      }
    },
    findGiftEffectiveTime: (value, that) => {
      for (let i = 0; i < that.GiftOptions.length; i++) {
        if (that.GiftOptions[i].mid === value) {
          if (that.GiftOptions[i].effectiveTime <= 0) {
            return "永久"
          } else {
            return that.unixTimeToDateTime(that.GiftOptions[i].effectiveTime)
          }
        }
      }
    },
  },
  methods: {
    GetGiftLogList() {
      this.loading = true
      GetGiftLogList(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.GetGiftLogList()
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000);
      const y = now.getFullYear();
      const m = now.getMonth() + 1;
      const d = now.getDate();
      return y + "-" + (m < 10 ? "0" + m : m) + "-" + (d < 10 ? "0" + d : d) + " " + now.toTimeString().substr(0, 8);
    },
  }
}
</script>

<style scoped>
.tab-container {
  margin: 30px;
}
</style>
