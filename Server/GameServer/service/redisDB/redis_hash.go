package redisDB

/*********************************************************
 *
 * Redis Hash
 *
 */

// 同时将多个 field-value (域-值)对设置到哈希表 key 中。会覆盖哈希表中已存在的字段，如果哈希表不存在，会创建一个空哈希表，并执行 HMSET 操作
func (this *RedisHelper) Hmset(args ...interface{}) (interface{}, error) {
	return this.Do("HMSET", args...)
}

// 获取所有给定字段的值 user na
func (this *RedisHelper) Hmget(args ...interface{}) (interface{}, error) {
	return this.Do("HMGET", args...)
}

// 将哈希表 key 中的字段 field 的值设为 value 。
func (this *RedisHelper) Hset(key, field, val interface{}) (interface{}, error) {
	return this.Do("HSET", key, field, val)
}

// 获取在哈希表中指定 key 的所有字段和值
func (this *RedisHelper) Hgetall(key interface{}) (interface{}, error) {
	return this.Do("HGETALL", key)
}

// 获取存储在哈希表中指定字段的值
func (this *RedisHelper) Hget(key, field interface{}) (interface{}, error) {
	return this.Do("HGET", key, field)
}

// 查看哈希表 key 中，指定的字段是否存在
func (this *RedisHelper) Hexists(key, field interface{}) (interface{}, error) {
	return this.Do("HEXISTS", key, field)
}

// 为哈希表 key 中的指定字段的整数值加上增量 increment
func (this *RedisHelper) Hincrby(key, field, amount interface{}) (interface{}, error) {
	return this.Do("HINCRBY", key, field, amount)
}

// 获取哈希表中字段的数量
func (this *RedisHelper) Hlen(key interface{}) (interface{}, error) {
	return this.Do("HLEN", key)
}

// 删除一个或多个哈希表字段
func (this *RedisHelper) Hdel(key interface{}, fields ...interface{}) (interface{}, error) {
	return this.Do("HDEL", append([]interface{}{key}, fields...)...)
}

// 获取哈希表中所有值
func (this *RedisHelper) Hvals(key interface{}) (interface{}, error) {
	return this.Do("HVALS", key)
}

// 为哈希表 key 中的指定字段的浮点数值加上增量 increment
func (this *RedisHelper) Hincrbyfloat(key, field, amount interface{}) (interface{}, error) {
	return this.Do("HINCRBYFLOAT", key, field, amount)
}

// 获取所有哈希表中的字段(KEY)
func (this *RedisHelper) Hkeys(key interface{}) (interface{}, error) {
	return this.Do("HKEYS", key)
}

// 只有在字段不存在时，设置哈希表字段的值
func (this *RedisHelper) Hsetnx(key, field, val interface{}) (interface{}, error) {
	return this.Do("HSETNX", key, field, val)
}
