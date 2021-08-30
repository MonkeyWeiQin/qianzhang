<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
<!--      <el-form-item label="用户名">-->
<!--        <el-input placeholder="用户名" v-model="listQuery.username"></el-input>-->
<!--      </el-form-item>-->
      <el-form-item>
<!--        <el-button type="primary" @click="getNoticeList()">查询</el-button>-->
        <el-button type="primary" @click="dialogFormCreate = true;resetFields()">添加公告</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="50" prop="mid" label="ID"/>
      <el-table-column width="300" prop="title" label="标题"></el-table-column>
      <el-table-column width="489" prop="content" label="内容">
        <el-row slot-scope="scope">
          <div v-html="scope.row.content"></div>
        </el-row>
      </el-table-column>
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
      <el-table-column width="150" label="操作">
        <el-row slot-scope="scope">
          <el-button type="text" size="mini" @click="dialogFormCreate = true;
           NoticeForm.startTime =scope.row.startTime * 1000;
           NoticeForm.endTime =scope.row.endTime * 1000;
           NoticeForm.title =scope.row.title;
           NoticeForm.content = scope.row.content;
           NoticeForm.mid = scope.row.mid;
           NoticeForm.type = scope.row.type;
          ">
            编辑
          </el-button>

        </el-row>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination background layout="prev, pager, next" @current-change="currentChange" :page-size="listQuery.limit"
                     :total="total"/>
    </div>
    <el-dialog title="新增/编辑公告" :visible.sync="dialogFormCreate">
      <el-form :model="NoticeForm" ref="NoticeForm" :rules="rules">
        <el-form-item label="开始时间" :label-width="formLabelWidth" prop="startTime" required>
          <el-date-picker type="datetime" v-model="NoticeForm.startTime" value-format="timestamp"></el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间" :label-width="formLabelWidth" prop="endTime" required>
          <el-date-picker type="datetime" v-model="NoticeForm.endTime" value-format="timestamp"></el-date-picker>
        </el-form-item>
        <el-form-item label="标题" :label-width="formLabelWidth" prop="title" required>
          <el-input v-model="NoticeForm.title" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="类型" :label-width="formLabelWidth" prop="type" hidden>
          <el-input v-model="NoticeForm.type" autocomplete="off" readonly ></el-input>
        </el-form-item>
        <el-form-item label="内容" :label-width="formLabelWidth" prop="content" >
          <el-input type="textarea" v-model="NoticeForm.content" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="default" class="dialog-footer" :label-width="formLabelWidth">
        <el-button @click="dialogFormCreate = false">取 消</el-button>
        <el-button v-if="NoticeForm.mid===''" type="primary" @click="CreateNotice">确 定</el-button>
        <el-button v-if="NoticeForm.mid!==''" type="primary" @click="UpdateNotice">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script>
// import Tinymce from '@/components/Tinymce'
import {getNoticeList, createNotice, updateNotice} from '@/api/notice'

export default {
  name: 'GameNoticeList',
  // components: {Tinymce},
  data() {
    return {
      formLabelWidth: '120px',
      dialogFormCreate: false,
      that: this,
      rules: {},
      list: [],
      total: 0,
      options: [],
      noticeType: 1,
      listQuery: {
        page: 1,
        limit: 5,
        type: 1,
      },
      NoticeForm: {
        content: "",
        type: 1,
        title: "",
        startTime: "",
        endTime: "",
        mid:""
      },
      DeleteGiftForm: {
        mid: "",
      }
    }
  },
  created() {
    this.getNoticeList()
  },
  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  methods: {
    getNoticeList() {
      this.loading = true
      getNoticeList(this.listQuery).then(response => {
        this.list = response.data.list
        this.loading = false
        this.total = response.data.total
      })
    },
    resetFields() {
      this.$nextTick(() => {
        this.$refs.NoticeForm.resetFields()
      })
    },
    currentChange(val) {
      this.listQuery.page = val
      this.getNoticeList()
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000);
      const y = now.getFullYear();
      const m = now.getMonth() + 1;
      const d = now.getDate();
      return y + "-" + (m < 10 ? "0" + m : m) + "-" + (d < 10 ? "0" + d : d) + " " + now.toTimeString().substr(0, 8);
    },
    CreateNotice() {
      this.NoticeForm.startTime /= 1000
      this.NoticeForm.endTime /= 1000
      createNotice(this.NoticeForm).then(response => {
        this.getNoticeList()
        this.dialogFormCreate = false
        this.resetFields()
      })
    },
    UpdateNotice() {
      this.NoticeForm.startTime /= 1000
      this.NoticeForm.endTime /= 1000
      updateNotice(this.NoticeForm).then(response => {
        this.getNoticeList()
        this.dialogFormCreate = false
        this.resetFields()
      })
    },
    DeleteGiftCode() {

    },
  }

}
</script>

<style scoped>
.tab-container {
  margin: 30px;
}
</style>
