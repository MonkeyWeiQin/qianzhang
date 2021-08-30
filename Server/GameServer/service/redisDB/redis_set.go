package redisDB

/*********************************************************
 *
 * Redis Set
 *
 */

// 返回给定集合的并集。不存在的集合 key 被视为空集
func (this *RedisHelper) Sunion(args ...interface{}) (interface{}, error) {
	return this.Do("SUNION", args...)
}

// 返回集合中元素的数量
func (this *RedisHelper) Scard(key interface{}) (interface{}, error) {
	return this.Do("SCARD", key)
}

// 返回集合中一个或多个随机数
func (this *RedisHelper) Srandmember(key interface{}, count ...interface{}) (interface{}, error) {
	if len(count) > 0 {
		return this.Do("SCARD", key, count[0])
	} else {
		return this.Do("SCARD", key)
	}
}

// 返回集合中的所有的成员。 不存在的集合 key 被视为空集合
func (this *RedisHelper) Smembers(key interface{}) (interface{}, error) {
	return this.Do("SMEMBERS", key)
}

// 返回给定所有给定集合的交集。 不存在的集合 key 被视为空集。 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)
func (this *RedisHelper) Sinter(args ...interface{}) (interface{}, error) {
	return this.Do("SINTER", args...)
}

// 移除集合中一个或多个成员
func (this *RedisHelper) Srem(args ...interface{}) (interface{}, error) {
	return this.Do("SREM", args...)
}

// 将指定成员 member 元素从 source 集合移动到 destination 集合
func (this *RedisHelper) Smove(source, destination, member interface{}) (interface{}, error) {
	return this.Do("SMOVE", source, destination, member)
}

// 向集合添加一个或多个成员
func (this *RedisHelper) Sadd(args ...interface{}) (interface{}, error) {
	return this.Do("SADD", args...)
}

//  判断 member 元素是否是集合 key 的成员
func (this *RedisHelper) Sismember(key, val interface{}) (interface{}, error) {
	return this.Do("SISMEMBER", key, val)
}

//  返回给定所有集合的差集并存储在 destination 中
func (this *RedisHelper) Sdiffstore(args ...interface{}) (interface{}, error) {
	return this.Do("SDIFFSTORE", args...)
}

//  返回给定集合之间的差集。不存在的集合 key 将视为空集
func (this *RedisHelper) Sdiff(args ...interface{}) (interface{}, error) {
	return this.Do("SDIFF", args...)
}

//  迭代集合中的元素
func (this *RedisHelper) Sscan(args ...interface{}) (interface{}, error) {
	return this.Do("SSCAN", args)
}

//  将给定集合之间的交集存储在指定的集合中。如果指定的集合已经存在，则将其覆盖。
func (this *RedisHelper) Sinterstore(args ...interface{}) (interface{}, error) {
	return this.Do("SINTERSTORE", args)
}

//  将给定集合的并集存储在指定的集合 destination 中
func (this *RedisHelper) Sunionstore(args ...interface{}) (interface{}, error) {
	return this.Do("SUNIONSTORE", args)
}

//  移除并返回集合中的一个随机元素
func (this *RedisHelper) Spop(key interface{}) (interface{}, error) {
	return this.Do("SPOP", key)
}
