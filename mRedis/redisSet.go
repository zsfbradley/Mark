package mRedis

import "markV5/mFace"

/*
经典使用场景:
1.标签系统:通过标签来识别不同兴趣的用户，来推荐带有相同标签的内容
2.生成随机数、抽奖:通过sPop、sRandMember来随机获取
3.社交需求:通过sAdd+sInter来识别推荐具有共同兴趣的用户
*/

/*
集合的特点:
1.无序，所以不能通过下标访问
2.不可重复，单个集合中不存在相同的元素
*/

// 向 key 的集合中添加 elements 元素
func (r *MRedis) SAdd(key string, elements ...interface{}) (interface{}, mFace.MError) {
	return r.Exec("SADD", key, elements...)
}

// 从 key 的集合中删除 elements 元素
func (r *MRedis) SRem(key string, elements ...interface{}) (interface{}, mFace.MError) {
	return r.Exec("SREM", key, elements...)
}

// 计算 key 的集合中元素个数
func (r *MRedis) SCard(key string) (interface{}, mFace.MError) {
	return r.Exec("SCARD", key)
}

// 判断 key 的集合中是否含有 element 元素
func (r *MRedis) SIsMember(key string, element interface{}) (interface{}, mFace.MError) {
	return r.Exec("SISMEMBER", key, element)
}

// 从 key 的集合中随机返回指定 count 个数的元素，count需大于等于1，元素不会被删除
func (r *MRedis) SRandMember(key string, count int) (interface{}, mFace.MError) {
	if count <= 0 {
		count = 1
	}
	return r.Exec("SRANDMEMBER", key, count)
}

// 从 key 的集合中随机弹出一个元素，元素会被删除
func (r *MRedis) SPop(key string) (interface{}, mFace.MError) {
	return r.Exec("SPOP", key)
}

// 获取 key 的集合所有元素，如果集合内元素较多，禁止在正式服使用
func (r *MRedis) SMembers(key string) (interface{}, mFace.MError) {
	return r.Exec("SMEMBERS", key)
}

// ------------------------ 集合间操作

// 求 keys 多个集合的交集
func (r *MRedis) SInter(keys ...string) (interface{}, mFace.MError) {
	param := make([]interface{}, 0)
	for _, key := range keys {
		param = append(param, key)
	}
	return r.baseExec("SINTER", param...)
}

// 求 keys 多个集合的交集并保存到指定的 destination key 集合中
func (r *MRedis) SInterStore(desKey string, keys ...string) (interface{}, mFace.MError) {
	param := make([]interface{}, 0)
	for _, key := range keys {
		param = append(param, key)
	}
	return r.Exec("SINTERSTORE", desKey, param...)
}

// 求 keys 多个集合的并集
func (r *MRedis) SUnion(keys ...string) (interface{}, mFace.MError) {
	param := make([]interface{}, 0)
	for _, key := range keys {
		param = append(param, key)
	}
	return r.baseExec("SUNION", param...)
}

// 求 keys 多个集合的并集并保存到指定的 destination key 集合中
func (r *MRedis) SUnionStore(desKey string, keys ...string) (interface{}, mFace.MError) {
	param := make([]interface{}, 0)
	for _, key := range keys {
		param = append(param, key)
	}
	return r.Exec("SUNIONSTORE", desKey, param...)
}

// 求 keys 多个集合的差集
// 结果的逻辑为，求第一个 key 的集合相对于后面的集合的内容的差集
func (r *MRedis) SDiff(keys ...string) (interface{}, mFace.MError) {
	param := make([]interface{}, 0)
	for _, key := range keys {
		param = append(param, key)
	}
	return r.baseExec("SDIFF", param...)
}

// 求 keys 多个集合的差集并保存到指定的 destination key 集合中
// 结果的逻辑为，求第一个 key 的集合相对于后面的集合的内容的差集
func (r *MRedis) SDiffStore(desKey string, keys ...string) (interface{}, mFace.MError) {
	param := make([]interface{}, 0)
	for _, key := range keys {
		param = append(param, key)
	}
	return r.Exec("SDIFFSTORE", desKey, param...)
}
