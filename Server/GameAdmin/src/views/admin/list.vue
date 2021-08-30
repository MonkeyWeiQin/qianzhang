<template>
  <div class="tab-container">
    <el-form v-model="listQuery" :inline="true" class="demo-form-inline">
<!--      <el-form-item label="用户名">-->
<!--        <el-input v-model="listQuery.username" placeholder="用户名" />-->
<!--      </el-form-item>-->
      <el-form-item>
<!--        <el-button type="primary" @click="getAdminList()">查询</el-button>-->
        <el-button type="primary" @click="dialogFormAddAdmin = true;AdminForm.add=true">新增管理员</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="100" prop="username" label="用户名" />
      <el-table-column width="80" prop="status" label="账号状态">
        <el-button slot-scope="scope" type="text">
          <p v-if="scope.row.status === 0 || scope.row.status < (new Date().getTime()) / 1000 ">
            <el-tag :type="'success'" disable-transitions>正常</el-tag>
          </p>
          <p v-else>
            <el-tag :type="'danger'" disable-transitions>被禁用（解封时间:）{{ scope.row.status|formatterTime(that) }}</el-tag>
          </p>
        </el-button>
      </el-table-column>
      <el-table-column width="170" label="最后登录时间">
        <template slot-scope="scope">
          <div v-if="scope.row.last_login_time > 0 ">
            <el-tag size="mini"> {{ scope.row.last_login_time|formatterTime(that) }}</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column width="200" label="操作">
        <el-row slot-scope="scope">
          <el-button
            type="text"
            size="mini"
            @click="dialogFormAddAdmin = true;AdminForm.add = false;AdminForm.username = scope.row.username;AdminForm.status = scope.row.status"
          >
            编辑
          </el-button>
          <el-button
            type="text"
            size="mini"
            @click="dialogFormUpdatePassword = true;UpdatePasswordForm.username = scope.row.username;"
          >
            修改密码
          </el-button>
          <el-button
            type="text"
            size="mini"
            @click="dialogFormRole = true;UpdateRoleForm.username=scope.row.username;UpdateRoleForm.role=scope.row.role;GetLevelMenuList()"
          >
            分配权限
          </el-button>
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
    <el-dialog title="添加/修改管理员" :visible.sync="dialogFormAddAdmin">
      <el-form ref="AdminForm" :model="AdminForm">
        <el-form-item label="账号" :label-width="formLabelWidth" required prop="username">
          <el-input v-model="AdminForm.username" autocomplete="off" />
        </el-form-item>
        <div v-if=" AdminForm.add === true ">
          <el-form-item label="密码" :label-width="formLabelWidth" required prop="password">
            <el-input v-model="AdminForm.password" autocomplete="off" />
          </el-form-item>
        </div>
        <el-form-item label="当前状态" :label-width="formLabelWidth" prop="status">
          <el-select v-model="AdminForm.status" placeholder="请选择">
            <el-option label="正常" :value="0" />
            <el-option label="停用" :value="-1" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormAddAdmin = false">取 消</el-button>
        <el-button type="primary" @click="AddAdmin">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="修改密码" :visible.sync="dialogFormUpdatePassword">
      <el-form ref="UpdatePasswordForm" :model="UpdatePasswordForm">
        <el-form-item label="账号" :label-width="formLabelWidth" required prop="username">
          <el-input v-model="UpdatePasswordForm.username" readonly disabled autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" :label-width="formLabelWidth" required prop="password">
          <el-input v-model="UpdatePasswordForm.password" autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormUpdatePassword = false">取 消</el-button>
        <el-button type="primary" @click="UpdatePassword">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="分配权限" :visible.sync="dialogFormRole">
      <el-tree
        ref="roleTree"
        :data="roleTree"
        show-checkbox
        default-expand-all
        node-key="mid"
        :default-checked-keys="UpdateRoleForm.role"
        highlight-current
        :props="defaultProps"
      />
      <div style="margin-top: 10px;">
        <div slot="default" class="dialog-footer">
          <el-button @click="dialogFormRole = false">取 消</el-button>
          <el-button type="primary" @click="UpdateRole">确 定</el-button>
        </div>
      </div>

    </el-dialog>
  </div>
</template>
<script>

import { AddAdmin, AdminList, UpdatePassword, UpdateRole } from '@/api/admin'
import { GetLevelMenuList } from '@/api/menu'

export default {
  name: 'AdminList',
  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  data() {
    return {
      formLabelWidth: '120px',
      dialogFormAddAdmin: false,
      dialogFormRole: false,
      dialogFormUpdatePassword: false,
      list: [],
      total: 0,
      that: this,
      roleTree: [],
      defaultProps: {
        children: 'children',
        label: 'name'
      },
      listQuery: {
        page: 1,
        limit: 5,
        username: ''
      },
      AdminForm: {
        username: '',
        password: '',
        status: 0,
        add: false
      },
      UpdatePasswordForm: {
        username: '',
        password: ''
      },
      UpdateRoleForm: {
        username: '',
        role: []
      }
    }
  },
  created() {
    this.getAdminList()
  },
  methods: {
    getAdminList() {
      this.loading = true
      AdminList(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.getAdminList()
    },
    GetLevelMenuList() {
      GetLevelMenuList().then((response) => {
        this.roleTree = response.data
      })
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },

    AddAdmin() {
      this.$refs.AdminForm.validate((valid) => {
        if (valid) {
          AddAdmin(this.AdminForm).then(() => {
            this.dialogFormAddAdmin = false
            this.getAdminList()
            this.$refs.AdminForm.resetFields()
          })
        }
      })
    },
    UpdatePassword() {
      this.$refs.UpdatePasswordForm.validate((valid) => {
        if (valid) {
          UpdatePassword(this.UpdatePasswordForm).then(() => {
            this.dialogFormUpdatePassword = false
            this.getAdminList()
            this.$refs.UpdatePasswordForm.resetFields()
          })
        }
      })
    },
    UpdateRole() {
      UpdateRole({ username: this.UpdateRoleForm.username, role: this.$refs.roleTree.getCheckedKeys() }).then(() => {
        this.dialogFormRole = false
        this.getAdminList()
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
