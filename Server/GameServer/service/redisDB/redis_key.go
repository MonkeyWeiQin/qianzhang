package redisDB

/*********************************************************
 *
 * Redis Key
 *
 */

// 同时将多个 field-value (域-值)对设置到哈希表 key 中。
func (this *RedisHelper) Type(key string) (interface{}, error) {
	return this.Do("TYPE", key)
}

// 设置 key 的过期时间 以毫秒计
func (this *RedisHelper) Pexpireat(key string, t int64) (interface{}, error) {
	return this.Do("PEXPIREAT", key, t)
}

// 修改 key 的名称
func (this *RedisHelper) Rename(key, newKey string) (interface{}, error) {
	return this.Do("RENAME", key, newKey)
}

// 移除给定 key 的过期时间，使得 key 永不过期
func (this *RedisHelper) Persist(key string) (interface{}, error) {
	return this.Do("PERSIST", key)
}

// 将当前数据库的 key 移动到给定的数据库 db 当中
func (this *RedisHelper) Move(key string, db int) (interface{}, error) {
	return this.Do("MOVE", key, db)
}

// 当前数据库中随机返回一个 key
func (this *RedisHelper) RandomKey() (interface{}, error) {
	return this.Do("RANDOMKEY")
}

// 序列化给定 key ，并返回被序列化的值
func (this *RedisHelper) Dump(key string) (interface{}, error) {
	return this.Do("DUMP", key)
}

// 以秒为单位返回 key 的剩余过期时间
func (this *RedisHelper) TTL(key string) (interface{}, error) {
	return this.Do("TTL", key)
}

// 设置 key 的过期时间, key 过期后将不再可用  单位 秒
func (this *RedisHelper) Expire(key string, t int64) (interface{}, error) {
	return this.Do("EXPIRE", key, t)
}

// 删除已存在的键。不存在的 key 会被忽略
func (this *RedisHelper) DEL(key string) (interface{}, error) {
	return this.Do("DEL", key)
}

// 以毫秒为单位返回 key 的剩余过期时间
func (this *RedisHelper) Pttl(key string) (interface{}, error) {
	return this.Do("PTTL", key)
}

// 修改key的名称 新名称不存在时修改，新的key名称已经存在于数据库则返回 0
func (this *RedisHelper) Renamenx(key, newKey string) (interface{}, error) {
	return this.Do("RENAMENX", key, newKey)
}

// 检查给定 key 是否存在
func (this *RedisHelper) Exists(key string) (interface{}, error) {
	return this.Do("EXISTS", key)
}

// 以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间。key 过期后将不再可用
func (this *RedisHelper) Expireat(key string, t int) (interface{}, error) {
	return this.Do("EXPIREAT", key, t)
}

// 查找所有符合给定模式 pattern 的 key
func (this *RedisHelper) Keys(filter string) (interface{}, error) {
	return this.Do("KEYS", filter)
}
