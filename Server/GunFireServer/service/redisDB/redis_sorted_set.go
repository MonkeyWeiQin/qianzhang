package redisDB

/*********************************************************
 *
 * Redis Sorted Set
 *
 */

// 命令返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序
func (this *RedisHelper) Zrevrank(key, member interface{}) (interface{}, error) {
	return this.Do("ZREVRANK", key, member)
}

// 在有序集合中计算指定字典区间内成员数量
func (this *RedisHelper) Zlexcount(key, min, max interface{}) (interface{}, error) {
	return this.Do("ZLEXCOUNT", key, key, min, max)
}

// 计算给定的一个或多个有序集的并集，并存储在新的 key 中
func (this *RedisHelper) Zunionstore(args ...interface{}) (interface{}, error) {
	return this.Do("ZUNIONSTORE", args...)
}

// 移除有序集合中给定的排名区间的所有成员
func (this *RedisHelper) Zremrangebyrank(key, start, stop interface{}) (interface{}, error) {
	return this.Do("ZREMRANGEBYRANK", key, start, stop)
}

// 获取有序集合的成员数
func (this *RedisHelper) Zcard(key interface{}) (interface{}, error) {
	return this.Do("ZCARD", key)
}

// 获取有序集合的成员数
func (this *RedisHelper) Zrem(args ...interface{}) (interface{}, error) {
	return this.Do("ZREM", args...)
}

// 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination
func (this *RedisHelper) Zinterstore(args ...interface{}) (interface{}, error) {
	return this.Do("ZINTERSTORE", args...)
}

// 返回有序集合中指定成员的索引
func (this *RedisHelper) Zrank(key, member interface{}) (interface{}, error) {
	return this.Do("ZRANK", key, member)
}

// 对有序集合中指定成员的分数加上增量 increment
func (this *RedisHelper) Zincrby(key, member, increment interface{}) (interface{}, error) {
	return this.Do("ZINCRBY", key, increment, member)
}

// 通过分数返回有序集合指定区间内的成员
func (this *RedisHelper) Zrangebyscore(args ...interface{}) (interface{}, error) {
	return this.Do("ZRANGEBYSCORE", args...)
}

// 迭代有序集合中的元素（包括元素成员和元素分值）
func (this *RedisHelper) Zscan(args ...interface{}) (interface{}, error) {
	return this.Do("ZSCAN", args...)
}

// 返回有序集中指定分数区间内的成员，分数从高到低排序
func (this *RedisHelper) Zrevrangebyscore(args ...interface{}) (interface{}, error) {
	return this.Do("ZREVRANGEBYSCORE", args...)
}

// 移除有序集合中给定的字典区间的所有成员
func (this *RedisHelper) Zremrangebylex(key, min, max interface{}) (interface{}, error) {
	return this.Do("ZREMRANGEBYLEX", key, min, max)
}

// 返回有序集中，指定区间内的成员
func (this *RedisHelper) Zrevrange(key, min, max interface{}) (interface{}, error) {
	return this.Do("ZREVRANGE", key, min, max)
}

// 用于计算有序集合中指定分数区间的成员数量
func (this *RedisHelper) ZADD(args ...interface{}) (interface{}, error) {
	return this.Do("ZADD", args...)
}

// 计算有序集合中指定分数区间的成员数量
func (this *RedisHelper) Zcount(key, min, max interface{}) (interface{}, error) {
	return this.Do("ZCOUNT", key, min, max)
}
