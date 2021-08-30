<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <el-form-item label="宝箱ID">
        <el-select v-model="listQuery.chestId" placeholder="请选择">
          <el-option label="全部" value=""></el-option>
          <el-option
            v-for="item in chestListOptions"
            :key="item.Id"
            :label="item.Id"
            :value="item.Id"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="用户ID">
        <el-input v-model="listQuery.uid" placeholder="用户ID"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="GetChestPurchaseList()">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="150" prop="chestId" label="宝箱ID">
      </el-table-column>
      <el-table-column width="100" prop="price" label="宝箱价格">
        <template slot-scope="scope">
          <span>{{ scope.row.price }}钻石</span>
        </template>
      </el-table-column>
      <el-table-column width="100" prop="uid" label="购买人">
        <template slot-scope="scope">
          <el-link type="primary" @click="GotoUser(scope.row.uid)">{{ scope.row.uid }}</el-link>
        </template>
      </el-table-column>
      <el-table-column width="100" prop="count" label="宝箱个数">
      </el-table-column>
      <el-table-column width="200" prop="uid" label="内容">
        <template slot-scope="scope">
          <div>{{ GetChestContent(scope.row.chestId, scope.row.chestItem) }}</div>
        </template>
      </el-table-column>
      <el-table-column width="170" label="购买时间">
        <template slot-scope="scope">
          <el-tag size="mini"> {{ scope.row.time|formatterTime(that) }}</el-tag>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-container">
      <el-pagination background layout="prev, pager, next" @current-change="currentChange" :page-size="listQuery.limit"
                     :total="total"
      />
    </div>
  </div>
</template>
<script>
import { GetChestList, GetChestPurchaseList, GetChestContent } from '@/api/shop'
import { commonConfig } from '@/common'

export default {
  name: 'BlindBoxList',
  data() {
    return {
      imgUrl: '',
      formLabelWidth: '120px',
      attachmentList: commonConfig.attachmentList,
      that: this,
      list: [],
      total: 0,
      chestContent: [],
      chestListOptions: [],
      listQuery: {
        page: 1,
        limit: 10,
        uid: '',
        chestId: ''
      },
      PriceForm: {
        type: '',
        unitPrice: '',
        repeatedlyPrice: ''
      },
      ItemForm: {
        type: '',
        item: []
      },
      SpecialItemForm: {
        type: '',
        item: []
      }
    }
  },
  created() {
    GetChestList().then(response => {
      this.chestListOptions = response.data
    })
    GetChestContent().then(response => {
      this.chestContent = response.data
    })
    this.GetChestPurchaseList()
  },
  filters: {
    formatterTime: (value, that) => {
      if (value) return that.unixTimeToDateTime(value)
      return ''
    }
  },
  methods: {
    currentChange(val) {
      this.listQuery.page = val
      this.GetChestPurchaseList()
    },
    getGoodsIcon(goodsId) {
      return require(`@/assets/shop/${goodsId}.png`)
    },
    GetChestPurchaseList() {
      this.loading = true
      GetChestPurchaseList(this.listQuery).then(response => {
        this.list = response.data.list
        this.total = response.data.total
        this.loading = false
      })
    },
    GetGoodsName(GoodsType) {
      let GoodsName = {}
      GoodsName[1] = '金币'
      GoodsName[2] = '钻石'
      GoodsName[3] = '体力'
      return GoodsName[GoodsType]
    },
    GotoUser(uid) {
      this.$router.push('/user/list?uid=' + uid)
    },
    GetChestContent(chestId, chestItem) {
      let textObj = {}
      let chestContent = this.chestContent[chestId]
      for (const key in chestContent) {
        chestItem.forEach(function(value) {
          if (chestContent.hasOwnProperty(key)) {
            if (value.itemId === chestContent[key].MatId) {
              if (!textObj.hasOwnProperty(chestContent[key].Des)) {
                textObj[chestContent[key].Des] = value.count
              } else {
                textObj[chestContent[key].Des] += value.count
              }
            }
          }
        })
      }
      let text = ""
      for (const item in textObj) {
        text += item +"*"+textObj[item] + " "
      }

      return text
    },
    unixTimeToDateTime: (time) => {
      const now = new Date(time * 1000)
      const y = now.getFullYear()
      const m = now.getMonth() + 1
      const d = now.getDate()
      return y + '-' + (m < 10 ? '0' + m : m) + '-' + (d < 10 ? '0' + d : d) + ' ' + now.toTimeString().substr(0, 8)
    }
  }
}
</script>

<style scoped>
.tab-container {
  margin: 30px;
}
</style>
