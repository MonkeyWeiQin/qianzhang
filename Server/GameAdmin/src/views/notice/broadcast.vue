<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <!--      <el-form-item label="广播类型">-->
      <!--        <el-select v-model="listQuery.position" placeholder="请选择">-->
      <!--          <el-option label="全部" value=""></el-option>-->
      <!--          <el-option-->
      <!--            v-for="item in itemTypeOptions"-->
      <!--            :key="item.value"-->
      <!--            :label="item.label"-->
      <!--            :value="item.value">-->
      <!--          </el-option>-->
      <!--        </el-select>-->
      <!--      </el-form-item>-->
      <el-form-item>
        <!--        <el-button type="primary" @click="getBroadcastList()">查询</el-button>-->
        <el-button type="primary" @click="dialogFormCreate = true;resetFields()">新增广播</el-button>
      </el-form-item>

    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="50" prop="mid" label="ID"/>
      <el-table-column width="180" prop="startTime" label="开始时间">
        <el-row slot-scope="scope">
          {{ scope.row.startTime|formatterTime(that) }}
        </el-row>
      </el-table-column>
      <el-table-column width="180" prop="endTime" label="结束时间">
        <el-row slot-scope="scope">
          {{ scope.row.endTime|formatterTime(that) }}
        </el-row>
      </el-table-column>
      <el-table-column width="100" prop="info" label="显示时间">
        <el-row slot-scope="scope">
          {{ scope.row.spacingTime|findSpacingTime(that) }}
        </el-row>
      </el-table-column>
      <!--      <el-table-column width="100" prop="position" label="展示位置">-->
      <!--        <el-row slot-scope="scope">-->
      <!--          {{ scope.row.position|findPosition(that) }}-->
      <!--        </el-row>-->
      <!--      </el-table-column>-->
      <el-table-column width="180" prop="color" label="显示颜色">
        <el-row slot-scope="scope">
          <div class="demo-color-box demo-color-box-other" :style="{background:scope.row.color}"><div class="value">{{scope.row.color}}</div></div>
        </el-row>
      </el-table-column>
      <el-table-column width="300" prop="content" label="广播内容"></el-table-column>

      <el-table-column width="150" label="操作">
        <el-row slot-scope="scope">
          <el-button type="text" size="mini" @click="dialogFormCreate = true;
           BroadcastForm.mid =scope.row.mid;
           BroadcastForm.content =scope.row.content;
           BroadcastForm.color =scope.row.color;
           BroadcastForm.startTime =scope.row.startTime*1000;
           BroadcastForm.endTime =scope.row.endTime*1000;
           BroadcastForm.spacingTime =scope.row.spacingTime;
          scope.row.effectiveTime>0?BroadcastForm.isForever = false:BroadcastForm.isForever=true;"
          >
            编辑
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
      <el-form :model="BroadcastForm" ref="BroadcastForm" :rules="rules">
        <el-form-item label="开始时间" :label-width="formLabelWidth" prop="startTime" required>
          <el-date-picker type="datetime" v-model="BroadcastForm.startTime" value-format="timestamp"></el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间" :label-width="formLabelWidth" prop="endTime" required>
          <el-date-picker type="datetime" v-model="BroadcastForm.endTime" value-format="timestamp"></el-date-picker>
        </el-form-item>
        <el-form-item label="展示颜色" :label-width="formLabelWidth" prop="color" required>
          <el-color-picker v-model="BroadcastForm.color"></el-color-picker>
        </el-form-item>
        <!--        <el-form-item label="是否重复发送" :label-width="formLabelWidth">-->
        <!--          <el-row :gutter="24">-->
        <!--            <el-col :span="8">-->
        <!--              <el-switch v-model="BroadcastForm.isForever" active-color="#13ce66" inactive-color="#ff4949"></el-switch>-->
        <!--            </el-col>-->
        <!--          </el-row>-->
        <!--        </el-form-item>-->
        <el-form-item label="显示时间(单位:秒)" :label-width="formLabelWidth" prop="spacingTime"
                      v-if="BroadcastForm.isForever!==false"
        >
          <el-row :gutter="24">
            <el-col :span="8">
              <el-input-number v-model="BroadcastForm.spacingTime" :min="5"></el-input-number>
            </el-col>
          </el-row>
        </el-form-item>
        <el-row :gutter="24">
          <el-col :span="10">
            <el-form-item label="内容" :label-width="formLabelWidth" prop="content" required>
              <el-input type="textarea" v-model="BroadcastForm.content" autocomplete="off"
                        :autosize="{ minRows: 3, maxRows: 6}"
              ></el-input>

            </el-form-item>
          </el-col>
        </el-row>

      </el-form>
      <div slot="default" class="dialog-footer" :label-width="formLabelWidth">
        <el-button @click="dialogFormCreate = false">取 消</el-button>
        <el-button v-if="BroadcastForm.mid===''" type="primary" @click="CreateBroadcast">确 定</el-button>
        <el-button v-if="BroadcastForm.mid!==''" type="primary" @click="UpdateBroadcast">确 定</el-button>
      </div>
    </el-dialog>

  </div>
</template>
<script>

import { getBroadcastList, createBroadcast, updateBroadcast } from '@/api/broadcast'

export default {
  name: 'BroadcastList',
  data() {
    return {
      formLabelWidth: '120px',
      dialogFormCreate: false,
      that: this,
      itemColorOptions: [
        {
          value: 1,
          label: '红色'
        },
        {
          value: 2,
          label: '蓝色'
        },
        {
          value: 3,
          label: '黄色'
        }
      ],
      list: [],
      total: 0,
      options: [],
      rules: {
      },
      listQuery: {
        page: 1,
        limit: 5,
        username: '',
        position: ''
      },
      BroadcastForm: {
        content: '',
        color: '',
        startTime: '',
        endTime: '',
        spacingTime: 0,
        mid: ''
      },
      DeleteBroadcastForm: {
        mid: ''
      }
    }
  },
  created() {
    this.getBroadcastList()
  },

  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    },
    findSpacingTime: (value, that) => {
      if (value === 0) {
        return '不重复'
      } else {
        return '显示' + value + '秒'
      }
    },
    findPosition: (value, that) => {
      console.log(value)
      for (let i = 0; i < that.itemTypeOptions.length; i++) {
        if (that.itemTypeOptions[i].value === value) {
          return that.itemTypeOptions[i].label
        }
      }
    },
    findColor: (value, that) => {
      for (let i = 0; i < that.itemColorOptions.length; i++) {
        if (that.itemColorOptions[i].value === value) {
          return that.itemColorOptions[i].label
        }
      }
    }
  },
  methods: {
    getBroadcastList() {
      this.loading = true
      getBroadcastList(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    resetFields() {
      this.$nextTick(() => {
        this.$refs.BroadcastForm.resetFields()
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.getBroadcastList()
    },

    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },
    CreateBroadcast() {
      this.BroadcastForm.startTime /= 1000
      this.BroadcastForm.endTime /= 1000

      createBroadcast(this.BroadcastForm).then(() => {
        this.getBroadcastList()
        this.dialogFormCreate = false
        this.resetFields()
      })
    },
    UpdateBroadcast() {
      this.BroadcastForm.startTime /= 1000
      this.BroadcastForm.endTime /= 1000

      updateBroadcast(this.BroadcastForm).then(() => {
        this.getBroadcastList()
        this.dialogFormCreate = false
        this.resetFields()
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
