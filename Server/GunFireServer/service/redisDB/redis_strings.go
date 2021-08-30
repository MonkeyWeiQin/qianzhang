package redisDB

/*********************************************************
 *
 * Redis String
 *
 */

// 指定的 key 不存在时，为 key 设置指定的值
func (this *RedisHelper) Setnx(key, val interface{}) (interface{}, error) {
	return this.Do("SETNX", key, val)
}

// 获取存储在指定 key 中字符串的子字符串。字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)
func (this *RedisHelper) Getrange(key, val interface{}) (interface{}, error) {
	return this.Do("GETRANGE", key, val)
}

// 同时设置一个或多个 key-value 对
func (this *RedisHelper) Mset(args ...interface{}) (interface{}, error) {
	return this.Do("MSET", args...)
}

// 为指定的 key 设置值及其过期时间。如果 key 已经存在， SETEX 命令将会替换旧的值
func (this *RedisHelper) Setex(key, val interface{}, timeout int64) (interface{}, error) {
	return this.Do("SETEX", key, timeout, val)
}

// 设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型
func (this *RedisHelper) SET(key, val interface{}) (interface{}, error) {
	return this.Do("SET", key, val)
}

// 获取指定 key 的值。如果 key 不存在，返回 nil 。如果key 储存的值不是字符串类型，返回一个错误
func (this *RedisHelper) GET(key interface{}) (interface{}, error) {
	return this.Do("GET", key)
}

// 对 key 所储存的字符串值，获取指定偏移量上的位(bit)
func (this *RedisHelper) Getbit(key interface{}, offset int64) (interface{}, error) {
	return this.Do("GETBIT", key, offset)
}

// 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)
func (this *RedisHelper) Setbit(key interface{}, offset int64) (interface{}, error) {
	return this.Do("GETBIT", key, offset)
}

// 将 key 中储存的数字值减一，若值不能表示为数字，那么返回一个错误，若key不存在会将值初始化为0，再进行减操作
func (this *RedisHelper) Decr(key interface{}) (interface{}, error) {
	return this.Do("DECR", key)
}

// 将 key 中储存的数字值减去给定的量值，若值不能表示为数字，那么返回一个错误，若key不存在会将值初始化为0，再进行减操作
func (this *RedisHelper) Dectby(key interface{}, amount int64) (interface{}, error) {
	return this.Do("DECRBY", key, amount)
}

// 获取指定 key 所储存的字符串值的长度。当 key 储存的不是字符串值时，返回一个错误
func (this *RedisHelper) Strlen(key interface{}) (interface{}, error) {
	return this.Do("STRLEN", key)
}

// 所有给定 key 都不存在时，同时设置一个或多个 key-value 对
func (this *RedisHelper) Msetnx(args ...interface{}) (interface{}, error) {
	return this.Do("MSETNX", args...)
}

// 将 key 中储存的数字值增一，若值不能表示为数字，那么返回一个错误，若key不存在会将值初始化为0，再进行增量操作
func (this *RedisHelper) Incr(key interface{}, amount int64) (interface{}, error) {
	return this.Do("INCR", key, amount)
}

// 将 key 中储存的数字值加上自定的增量，若值不能表示为数字，那么返回一个错误，若key不存在会将值初始化为0，再进行增量操作
func (this *RedisHelper) Incrby(key interface{}, amount int64) (interface{}, error) {
	return this.Do("INCRBY", key, amount)
}

// 为 key 中所储存的值加上指定的浮点数增量值，key不存在将初始化为0 再进行操作
func (this *RedisHelper) Incrbyfloat(key interface{}, amount float64) (interface{}, error) {
	return this.Do("INCRBYFLOAT", key, amount)
}

// 指定的字符串覆盖给定 key 所储存的字符串值，覆盖的位置从偏移量 offset 开始
func (this *RedisHelper) Setrange(key, val interface{}, offset int64) (interface{}, error) {
	return this.Do("SETRANGE", key, offset, val)
}

// 以毫秒为单位设置 key 的生存时间
func (this *RedisHelper) Psetex(key, val interface{}, t int64) (interface{}, error) {
	return this.Do("PSETEX", key, val, t)
}

// 为指定的 key 追加值
func (this *RedisHelper) Append(key, val interface{}) (interface{}, error) {
	return this.Do("APPEND", key, val)
}

// 设置指定 key 的值，并返回 key 旧的值
func (this *RedisHelper) Getset(key, val interface{}) (interface{}, error) {
	return this.Do("GETSET", key, val)
}

// 返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil
func (this *RedisHelper) Mget(key, val interface{}) (interface{}, error) {
	return this.Do("MGET", key, val)
}
