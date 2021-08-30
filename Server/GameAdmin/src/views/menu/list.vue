<template>
  <div class="tab-container">
    <el-form v-model="listQuery" :inline="true" class="demo-form-inline">
<!--      <el-form-item label="用户名">-->
<!--        <el-input v-model="listQuery.username" placeholder="用户名" />-->
<!--      </el-form-item>-->
      <el-form-item>
        <el-button type="primary" @click="GetMenuList()">刷新</el-button>
        <el-button type="primary" @click="dialogFormCreate = true;resetFields(); MenuForm.mid='';GetLevelMenuList()">新增权限</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="50" prop="mid" label="ID" />
      <el-table-column width="300" prop="name" label="权限名称" />
      <el-table-column width="80" prop="type" label="权限类型">
        <div slot-scope="scope">{{ scope.row.type === 1 ? '目录' : '接口' }}</div>
      </el-table-column>
      <el-table-column width="300" prop="path" label="权限路由" />
<!--      <el-table-column width="80" prop="status" label="当前状态" />-->
<!--      <el-table-column width="150" prop="component" label="前端组件" />-->
<!--      <el-table-column width="80" prop="icon" label="icon" />-->

      <el-table-column width="150" label="操作">
        <el-row slot-scope="scope">
          <el-button
            type="text"
            size="mini"
            @click="dialogFormCreate = true;
                    MenuForm.mid =scope.row.mid;
                    MenuForm.pmid =scope.row.pmid;
                    MenuForm.type =scope.row.type;
                    MenuForm.name =scope.row.name;
                    MenuForm.path =scope.row.path;
                    MenuForm.component =scope.row.component;
                    MenuForm.title =scope.row.title;
                    MenuForm.icon =scope.row.icon;
                    MenuForm.status =scope.row.status;"
          >
            编辑
          </el-button>
          <el-button
            type="text"
            size="mini"
            @click="centerDialogVisible = true;DeleteMenuForm.mid=scope.row.mid"
          >
            删除
          </el-button>
        </el-row>
      </el-table-column>
    </el-table>

    <el-dialog title="新增/编辑权限" :visible.sync="dialogFormCreate">
      <el-form ref="MenuForm" :model="MenuForm" :rules="rules">
        <el-form-item label="上级节点" :label-width="formLabelWidth" required prop="pmid">
          <el-cascader
            v-model="MenuForm.pmid"
            :options="options"
            :props="{ multiple: false, checkStrictly: true ,value:'mid',label:'name'}"
          />
        </el-form-item>
        <el-form-item label="权限类型" :label-width="formLabelWidth" required prop="type">
          <el-radio-group v-model="MenuForm.type">
            <el-radio :label="1">菜单</el-radio>
            <el-radio :label="2">接口</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="权限名称" :label-width="formLabelWidth" required prop="name">
          <el-input v-model="MenuForm.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="权限路径" :label-width="formLabelWidth" required prop="path">
          <el-input v-model="MenuForm.path" autocomplete="off" />
        </el-form-item>
        <!--        <el-form-item label="前端组件key" :label-width="formLabelWidth" prop="component">-->
        <!--          <el-input v-model="MenuForm.component" autocomplete="off"></el-input>-->
        <!--        </el-form-item>-->
        <!--        <el-form-item label="图标" :label-width="formLabelWidth" prop="icon">-->
        <!--          <el-input v-model="MenuForm.icon" autocomplete="off"></el-input>-->
        <!--        </el-form-item>-->
        <!--        <el-form-item label="是否禁用" :label-width="formLabelWidth" prop="status">-->
        <!--          <el-switch-->
        <!--            v-model="MenuForm.status"-->
        <!--            active-color="#13ce66"-->
        <!--            inactive-color="#ff4949">-->
        <!--          </el-switch>-->
        <!--        </el-form-item>-->
      </el-form>
      <div slot="default" class="dialog-footer">
        <el-button @click="dialogFormCreate = false">取 消</el-button>
        <span v-if="MenuForm.mid!==''">
          <el-button type="primary" @click="UpdateMenu">确 定</el-button>
        </span>
        <span v-else>
          <el-button type="primary" @click="CreateMenu">确 定</el-button>
        </span>
      </div>
    </el-dialog>
    <el-dialog
      title="提示"
      :visible.sync="centerDialogVisible"
      width="30%"
      center
    >
      <span>确定删除该权限</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="centerDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="centerDialogVisible = false;DeleteMenu()">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>

import { GetMenuList, CreateMenu, UpdateMenu, GetLevelMenuList, DeleteMenu } from '@/api/menu'

export default {
  name: 'MenuList',

  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  data() {
    return {
      formLabelWidth: '120px',
      dialogFormCreate: false,
      dialogFormUpdate: false,
      centerDialogVisible: false,
      list: [],
      that: this,
      options: [],
      rules: {
        pmid: [
          { required: true, message: '请选择上级菜单', trigger: 'change' }
        ],
        type: [
          { required: true, message: '请选择权限类型', trigger: 'change' }
        ]
      },
      listQuery: {
        page: 1,
        limit: 10
      },
      MenuForm: {
        mid: '',
        pmid: '',
        type: '',
        name: '',
        path: '',
        component: '',
        title: '',
        icon: '',
        status: false
      },
      DeleteMenuForm: {
        mid: ''
      }
    }
  },
  created() {
    this.GetLevelMenuList()
    this.GetMenuList()
  },
  methods: {
    GetMenuList() {
      this.loading = true
      GetMenuList(this.listQuery).then(response => {
        this.list = response.data
        this.loading = false
      })
    },
    resetFields() {
      this.$nextTick(() => {
        this.$refs.MenuForm.resetFields()
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.GetMenuList()
    },

    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    },
    GetLevelMenuList() {
      GetLevelMenuList({ 'type': 1 }).then((response) => {
        this.options = response.data
      })
    },
    CreateMenu() {
      this.$refs.MenuForm.validate((valid) => {
        if (valid) {
          this.MenuForm.pmid = Number(this.MenuForm.pmid[this.MenuForm.pmid.length - 1])
          CreateMenu(this.MenuForm).then(() => {
            this.dialogFormCreate = false
            this.GetMenuList()
            this.resetFields()
          })
        } else {
          return false
        }
      })
    },
    UpdateMenu() {
      this.$refs.MenuForm.validate((valid) => {
        if (valid) {
          this.MenuForm.pmid = Number(this.MenuForm.pmid[this.MenuForm.pmid.length - 1])
          UpdateMenu(this.MenuForm).then(() => {
            this.dialogFormCreate = false
            this.GetMenuList()
            this.resetFields()
          })
        } else {
          return false
        }
      })
    },
    DeleteMenu() {
      DeleteMenu(this.DeleteMenuForm).then(() => {
        this.GetMenuList()
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
