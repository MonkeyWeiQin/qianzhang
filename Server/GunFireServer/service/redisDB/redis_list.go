package redisDB

/*********************************************************
 *
 * Redis List
 *
 */

// 通过索引获取列表中的元素。你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
func (this *RedisHelper) Lindex(key, index interface{}) (interface{}, error) {
	return this.Do("LINDEX", key, index)
}

// 将一个或多个值插入到列表的尾部(最右边)
func (this *RedisHelper) Rpush(args ...interface{}) (interface{}, error) {
	return this.Do("LINDEX", args...)
}

// 返回列表中指定区间内的元素，区间以偏移量 START 和 END 指定。
// 其中 0 表示列表的第一个元素， 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (this *RedisHelper) Lrange(key, start, end interface{}) (interface{}, error) {
	return this.Do("LRANGE", key, start, end)
}

// 移除列表的最后一个元素，并将该元素添加到另一个列表并返回。
func (this *RedisHelper) Rpoplpush(sourceListKey, destinationListKey interface{}) (interface{}, error) {
	return this.Do("RPOPLPUSH", sourceListKey, destinationListKey)
}

// 移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
func (this *RedisHelper) Blpop(args ...interface{}) (interface{}, error) {
	return this.Do("BLPOP", args...)
}

// 移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func (this *RedisHelper) Brpop(args ...interface{}) (interface{}, error) {
	return this.Do("BRPOP", args...)
}

// 列表中弹出一个值，将弹出的元素插入到另外一个列表中并返回它； 如果列表没有元素会阻塞列表直到等待超
func (this *RedisHelper) Brpoplpush(listKey, anotherListKey, t interface{}) (interface{}, error) {
	return this.Do("BRPOPLPUSH", listKey, anotherListKey, t)
}

// 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素。
//	COUNT 的值可以是以下几种：
//		count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
//		count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
// 		count = 0 : 移除表中所有与 VALUE 相等的值。
func (this *RedisHelper) Lrem(key, val, count interface{}) (interface{}, error) {
	return this.Do("LREM", key, count, val)
}

// 命令用于返回列表的长度。 如果列表 key 不存在，则 key 被解释为一个空列表，返回 0 。 如果 key 不是列表类型，返回一个错误。
func (this *RedisHelper) Llen(key interface{}) (interface{}, error) {
	return this.Do("LLEN", key)
}

//  对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
func (this *RedisHelper) Ltrim(key, start, stop interface{}) (interface{}, error) {
	return this.Do("LTRIM", key, start, stop)
}

//  移除并返回列表的第一个元素
func (this *RedisHelper) Lpop(key interface{}) (interface{}, error) {
	return this.Do("LPOP", key)
}

//  将一个或多个值插入到已存在的列表头部，列表不存在时操作无效。
func (this *RedisHelper) Lpushx(args ...interface{}) (interface{}, error) {
	return this.Do("LPUSHX", args...)
}

//  在列表的元素前或者后插入元素。 当指定元素不存在于列表中时，不执行任何操作。 当列表不存在时，被视为空列表，不执行任何操作。 如果 key 不是列表类型，返回一个错误。
func (this *RedisHelper) Linsert(key, option, val, newVal interface{}) (interface{}, error) {
	return this.Do("LINSERT", key, option, val, newVal)
}

//  移除并获取列表最后一个元素
func (this *RedisHelper) Rpop(key interface{}) (interface{}, error) {
	return this.Do("RPOP", key)
}

//  通过索引设置列表元素的值
func (this *RedisHelper) Lset(key, index, val interface{}) (interface{}, error) {
	return this.Do("LSET", key, index, val)
}

//  将一个或多个值插入到列表头部 如果key不存在则创建
func (this *RedisHelper) Lpush(args ...interface{}) (interface{}, error) {
	return this.Do("LPUSH", args...)
}

//  将一个或多个值插入到已存在的列表尾部(最右边)。如果列表不存在，操作无效
func (this *RedisHelper) Rpushx(args ...interface{}) (interface{}, error) {
	return this.Do("RPUSHX", args...)
}
