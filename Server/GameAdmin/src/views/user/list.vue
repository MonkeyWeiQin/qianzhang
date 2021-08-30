<template>
  <div class="tab-container">
    <el-form v-model="listQuery" :inline="true" class="demo-form-inline">
      <el-form-item label="用户ID">
        <el-input v-model="listQuery.uid" placeholder="用户ID"/>
      </el-form-item>
      <el-form-item label="用户名">
        <el-input v-model="listQuery.username" placeholder="用户名"/>
      </el-form-item>
      <el-form-item label="注册时间">
        <el-date-picker
          v-model="listQuery.register_time"
          type="datetimerange"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          range-separator="至"
          style="width: 100%;"
          value-format="timestamp"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="getUserList()">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="70" prop="uid" fixed="left" label="用户ID"/>
      <el-table-column width="100" prop="username" label="用户名"/>
      <el-table-column width="100" prop="mobile" label="用户手机号"/>
      <el-table-column width="110" prop="level" label="当前账号等级"/>
      <el-table-column width="110" prop="exp" label="当前账号经验"/>
      <el-table-column width="100" prop="diamond" label="钻石">
        <el-button
          slot-scope="scope"
          type="text"
          @click="dialogFormDiamond = true;diamondForm.uid = scope.row.uid;diamondForm.diamond = scope.row.diamond"
        >
          {{ scope.row.diamond }}
        </el-button>
      </el-table-column>
      <el-table-column width="100" prop="gold" label="金币">
        <el-button
          slot-scope="scope"
          type="text"
          @click="dialogFormGold = true;GoldForm.uid = scope.row.uid;GoldForm.gold = scope.row.gold"
        >
          {{ scope.row.gold }}
        </el-button>
      </el-table-column>
      <el-table-column width="70" prop="strength" label="体力值">
        <el-button
          slot-scope="scope"
          type="text"
          @click="dialogFormStrength = true;StrengthForm.uid = scope.row.uid;StrengthForm.strength = scope.row.strength"
        >
          {{ scope.row.strength }}
        </el-button>
      </el-table-column>
      <el-table-column width="80" prop="vip" label="VIP等级"/>
      <el-table-column width="80" prop="status" label="账号状态">
        <el-button
          slot-scope="scope"
          type="text"
          @click="dialogFormStatus = true;StatusForm.uid = scope.row.uid;StatusForm.status = scope.row.status"
        >
          <p v-if="scope.row.status === 0 || scope.row.status < (new Date().getTime()) / 1000 ">
            <el-tag :type="'success'" disable-transitions>正常</el-tag>
          </p>
          <p v-else>
            <el-tag :type="'danger'" disable-transitions>被禁用（解封时间:）{{ scope.row.status|formatterTime(that) }}</el-tag>
          </p>
        </el-button>

      </el-table-column>
      <el-table-column width="90" prop="vip_exp" label="VIP经验值"/>
      <el-table-column width="100" prop="vip_end_time" label="VIP到期时间"/>
      <el-table-column width="170" label="注册时间/最后登录时间">
        <template slot-scope="scope">
          <el-tag size="mini"> {{ scope.row.registerTime|formatterTime(that) }}</el-tag>
          <el-tag size="mini"> {{ scope.row.loginTime|formatterTime(that) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column width="150" label="操作">
        <el-row>
<!--          <el-button type="text" size="mini" >编辑</el-button>-->
        </el-row>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination
        background
        layout="prev, pager, next"
        :page-size="listQuery.limit"
        :total="total"
        @current-change="currentChange"
      />
    </div>
    <el-dialog title="调整钻石" :visible.sync="dialogFormDiamond" prop="uid">
      <el-form ref="diamondForm" :model="diamondForm">
        <el-form-item label="用户ID" :label-width="formLabelWidth">
          <el-input v-model="diamondForm.uid" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="调整方式" :label-width="formLabelWidth" prop="type">
          <el-select v-model="diamondForm.type" placeholder="请选择">
            <el-option label="增加" value="add"/>
            <el-option label="减少" value="reduce"/>
          </el-select>
        </el-form-item>
        <el-form-item label="当前数量" :label-width="formLabelWidth" prop="diamond">
          <el-input v-model="diamondForm.diamond" autocomplete="off" readonly disabled/>
        </el-form-item>
        <el-form-item label="数量" :label-width="formLabelWidth" prop="number">
          <el-input v-model="diamondForm.number" autocomplete="off"/>
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormDiamond = false">取 消</el-button>
        <el-button type="primary" @click="modifyDiamond">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="调整金币" :visible.sync="dialogFormGold">
      <el-form ref="GoldForm" :model="GoldForm">
        <el-form-item label="用户ID" :label-width="formLabelWidth" prop="uid">
          <el-input v-model="GoldForm.uid" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="调整方式" :label-width="formLabelWidth" prop="type">
          <el-select v-model="GoldForm.type" placeholder="请选择">
            <el-option label="增加" value="add"/>
            <el-option label="减少" value="reduce"/>
          </el-select>
        </el-form-item>
        <el-form-item label="当前数量" :label-width="formLabelWidth" prop="gold">
          <el-input v-model="GoldForm.gold" autocomplete="off" readonly disabled/>
        </el-form-item>
        <el-form-item label="数量" :label-width="formLabelWidth" prop="number">
          <el-input v-model="GoldForm.number" autocomplete="off"/>
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormGold = false">取 消</el-button>
        <el-button type="primary" @click="modifyGold">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="调整体力" :visible.sync="dialogFormStrength">
      <el-form ref="StrengthForm" :model="StrengthForm">
        <el-form-item label="用户ID" :label-width="formLabelWidth" prop="uid">
          <el-input v-model="StrengthForm.uid" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="调整方式" :label-width="formLabelWidth" prop="type">
          <el-select v-model="StrengthForm.type" placeholder="请选择">
            <el-option label="增加" value="add"/>
            <el-option label="减少" value="reduce"/>
          </el-select>
        </el-form-item>
        <el-form-item label="当前数量" :label-width="formLabelWidth" prop="strength">
          <el-input v-model="StrengthForm.strength" autocomplete="off" readonly disabled/>
        </el-form-item>
        <el-form-item label="数量" :label-width="formLabelWidth" prop="number">
          <el-input v-model="StrengthForm.number" autocomplete="off"/>
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormGold = false">取 消</el-button>
        <el-button type="primary" @click="modifyStrength">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="账号封禁/解封" :visible.sync="dialogFormStatus">
      <el-form ref="StatusForm" :model="StatusForm">
        <el-form-item label="用户ID" :label-width="formLabelWidth" prop="uid">
          <el-input v-model="StatusForm.uid" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="当前状态" :label-width="formLabelWidth" prop="status">
          <el-select v-model="StatusForm.status" placeholder="请选择">
            <el-option label="正常" :value="0"/>
            <el-option label="封禁" :value="-1"/>
          </el-select>
        </el-form-item>
        <div v-if="StatusForm.status !==0">
          <el-form-item label="解封时间" :label-width="formLabelWidth" prop="time">
            <el-date-picker
              v-model="StatusForm.time"
              type="date"
              placeholder="默认永久封号"
              :default-value="new Date()"
              style="width: 100%;"
              value-format="timestamp"
            />
          </el-form-item>
        </div>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormStatus = false">取 消</el-button>
        <el-button type="primary" @click="modifyStatus">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script>

import { modifyDiamond, modifyGold, modifyStatus, modifyStrength, userList } from '@/api/user-management'

export default {
  name: 'UserManagementList',

  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  data() {
    return {
      formLabelWidth: '120px',
      dialogFormDiamond: false,
      dialogFormGold: false,
      dialogFormStrength: false,
      dialogFormStatus: false,
      list: [],
      total: 0,
      that: this,
      listQuery: {
        page: 1,
        limit: 5,
        uid: '',
        username: '',
        register_time: [],
        register_end_time: '',
        register_start_time: ''
      },
      diamondForm: {
        uid: 0,
        diamond: 0,
        number: 0,
        type: ''
      },
      GoldForm: {
        uid: 0,
        gold: 0,
        number: 0,
        type: ''
      },
      StrengthForm: {
        uid: 0,
        strength: 0,
        number: 0,
        type: ''
      },
      StatusForm: {
        uid: 0,
        status: '',
        time: 0
      }
    }
  },
  created() {
    this.listQuery.uid = this.$route.query.uid
    this.getUserList()
  },

  methods: {
    getUserList() {
      this.loading = true
      if (this.listQuery.register_time && this.listQuery.register_time.length > 0) {
        this.listQuery.register_start_time = this.listQuery.register_time[0] / 1000
        this.listQuery.register_end_time = this.listQuery.register_time[1] / 1000
      } else {
        this.listQuery.register_start_time = this.listQuery.register_end_time = 0
      }
      userList(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.getUserList()
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },
    modifyDiamond() {
      this.diamondForm.number = Number(this.diamondForm.number)
      if (this.diamondForm.type === '' || this.diamondForm.number === '') {
        this.$message({ message: '请输入完整', type: 'warning' })
        return
      }
      if (this.diamondForm.type === 'reduce') {
        if (this.diamondForm.diamond < this.diamondForm.number) {
          this.$message({ message: '减少的钻石数量小于当前拥有的钻石数量', type: 'warning' })
          return
        }
      }
      modifyDiamond(this.diamondForm).then(() => {
        this.dialogFormDiamond = false
        this.getUserList()
        this.$refs.diamondForm.resetFields()
      })
    },
    modifyGold() {
      this.GoldForm.number = Number(this.GoldForm.number)
      if (this.GoldForm.type === '' || this.GoldForm.number === '') {
        this.$message({ message: '请输入完整', type: 'warning' })
        return
      }
      if (this.GoldForm.type === 'reduce') {
        if (this.GoldForm.gold < this.GoldForm.number) {
          this.$message({ message: '减少的金币数量小于当前拥有的金币数量', type: 'warning' })
          return
        }
      }
      modifyGold(this.GoldForm).then(() => {
        this.dialogFormGold = false
        this.getUserList()
        this.$refs.GoldForm.resetFields()
      })
    },
    modifyStrength() {
      this.StrengthForm.number = Number(this.StrengthForm.number)
      if (this.StrengthForm.type === '' || this.StrengthForm.number === '') {
        this.$message({ message: '请输入完整', type: 'warning' })
        return
      }
      if (this.StrengthForm.type === 'reduce') {
        if (this.StrengthForm.gold < this.StrengthForm.number) {
          this.$message({ message: '减少的体力数量小于当前拥有的体力数量', type: 'warning' })
          return
        }
      }
      modifyStrength(this.StrengthForm).then(() => {
        this.dialogFormStrength = false
        this.getUserList()
        this.$refs.StrengthForm.resetFields()
      })
    },
    modifyStatus() {
      if (this.StatusForm.status !== 0 && this.StatusForm.time !== 0) {
        this.StatusForm.time = this.StatusForm.time / 1000
      }
      modifyStatus(this.StatusForm).then(() => {
        this.dialogFormStatus = false
        this.getUserList()
        this.$refs.StatusForm.resetFields()
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
