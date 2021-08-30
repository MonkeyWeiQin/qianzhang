<template>
  <div class="tab-container">
    <el-form :inline="true" v-model="listQuery" class="demo-form-inline">
      <el-form-item label="商品Id">
        <el-select v-model="listQuery.goodsId" placeholder="请选择">
          <el-option label="全部" value=""></el-option>
          <el-option
            v-for="item in goodsListOptions"
            :key="item.Id"
            :label="item.Id"
            :value="item.Id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="用户ID">
        <el-input v-model="listQuery.uid" placeholder="用户ID"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="GetPurchaseList()">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" border style="width: 100%">
      <el-table-column width="100" prop="goodsId" label="商品icon">
        <template slot-scope="scope">
          <el-image class="table-td-thumb" :src="getGoodsIcon(scope.row.goodsId)"></el-image>
        </template>
      </el-table-column>
      <el-table-column width="100" prop="price" label="商品价格">
        <template slot-scope="scope">
          <span>{{ scope.row.price }}<span v-if="scope.row.goodsType === 2">元</span><span v-else>钻石</span></span>
        </template>
      </el-table-column>
      <el-table-column width="100" prop="uid" label="购买人">
        <template slot-scope="scope">
          <el-link type="primary" @click="GotoUser(scope.row.uid)">{{ scope.row.uid }}</el-link>
        </template>
      </el-table-column>
      <el-table-column width="200" prop="count" label="获得商品数量">
        <template slot-scope="scope">
          <span>{{ scope.row.count }}{{ GetGoodsName(scope.row.goodsType) }}</span>
          <span v-if="scope.row.presentation > 0">+<span>{{
              scope.row.presentation
            }}{{ GetGoodsName(scope.row.goodsType) }}
              <span></span>
          </span></span>
        </template>
      </el-table-column>
      <el-table-column width="200" prop="uid" label="描述">
        <template slot-scope="scope">
          <span v-if="scope.row.goodsType === 2">花费<span>{{ scope.row.price }}</span>元,获得{{
              scope.row.count
            }}{{ GetGoodsName(scope.row.goodsType) }}
            <span v-if="scope.row.presentation > 0">, 赠送<span>{{
                scope.row.presentation
              }}{{ GetGoodsName(scope.row.goodsType) }}</span></span>
          </span>
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
import { GetGoodsList, GetPurchaseList } from '@/api/shop'
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
      goodsListOptions: [],
      listQuery: {
        page: 1,
        limit: 6,
        uid: '',
        goodsId:''
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
    GetGoodsList().then(response=>{
      this.goodsListOptions = response.data
    })
    this.GetPurchaseList()
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
      this.GetPurchaseList()
    },
    getGoodsIcon(goodsId) {
      return require(`@/assets/shop/${goodsId}.png`)
    },
    GetPurchaseList() {
      this.loading = true
      GetPurchaseList(this.listQuery).then(response => {
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
