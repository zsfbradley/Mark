package mRedis

import (
	"markV5/mError"
	"markV5/mFace"
)

/*
经典使用场景:
1.消息队列:使用lpush+brpop命令的组合，可以实现左进消息，多客户端右阻塞式弹出的消息队列
2.文章列表:文章需要分页展示，列表不但有序，且支持按照索引范围获取元素
*/

/*
列表特点:
1.元素有序，意味着可以通过下标来获取指定的某个或某范围的元素
2.元素可重复
*/

type LInsert_Type string

const (
	LInsert_Before LInsert_Type = "before"
	LInsert_After  LInsert_Type = "after"
)

// 从右边添加元素或元素组
func (r *MRedis) RPush(key string, values ...interface{}) (interface{}, mFace.MError) {
	return r.Exec("RPUSH", key, values...)
}

// 从左边边添加元素或元素组
func (r *MRedis) LPush(key string, values ...interface{}) (interface{}, mFace.MError) {
	return r.Exec("LPUSH", key, values...)
}

// 在 target element 元素前或后插入指定元素 value
func (r *MRedis) LInsert(key string, insertType LInsert_Type, targetElement interface{}, value interface{}) (interface{}, mFace.MError) {
	if insertType != LInsert_Before && insertType != LInsert_After {
		return nil, mError.ErrorParam
	}
	return r.Exec("LINSERT", key, insertType, targetElement, value)
}

// 查找 key 从 start - end 范围内的元素，且包含 end
// 从左至右为 0 至 N-1
// 从右至左为 -1 至 -N
// start 必须小于 end
func (r *MRedis) LRange(key string, start, end int) (interface{}, mFace.MError) {
	return r.Exec("LRANGE", key, start, end)
}

// 获取指定 index 下标元素的值 ，index 从0开始
func (r *MRedis) LIndex(key string, index int) (interface{}, mFace.MError) {
	return r.Exec("LINDEX", key, index)
}

// 获取 key 列表的长度
func (r *MRedis) LLen(key string) (interface{}, mFace.MError) {
	return r.Exec("LLEN", key)
}

// 删除并返回 key 列表中最左侧的值
func (r *MRedis) LLPop(key string) (interface{}, mFace.MError) {
	return r.Exec("LPOP", key)
}

// 删除并返回 key 列表中最右侧的值
func (r *MRedis) LRPop(key string) (interface{}, mFace.MError) {
	return r.Exec("RPOP", key)
}

// 在 key 列表中查找等于 value 的元素，根据 count 的情况进行删除
// count > 0 ：从左至右删除最多 count 个元素
// count < 0 ：从右至左删除最多 count 的绝对值个元素
// count = 0 ：删除全部元素
func (r *MRedis) LRem(key string, count int, value interface{}) (interface{}, mFace.MError) {
	return r.Exec("LREM", key, count, value)
}

// 按照 start - end 裁剪 key 列表，保留 end 下标元素
func (r *MRedis) LTrim(key string, start, end int) (interface{}, mFace.MError) {
	return r.Exec("LTRIM", key, start, end)
}

// 修改 key 列表指定 index 下标元素的值为 newValue
func (r *MRedis) LSet(key string, index int, newValue interface{}) (interface{}, mFace.MError) {
	return r.Exec("LSET", key, index, newValue)
}

// 阻塞左弹出
// 扫描顺序由 keys 的顺序决定
// timeout 决定了阻塞时间，如果为0且 keys 皆为空列表，则一直阻塞
func (r *MRedis) LBLPop(timeout int, keys ...string) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for _, key := range keys {
		params = append(params, key)
	}
	params = append(params, timeout)
	return r.baseExec("BLPOP", params...)
}
