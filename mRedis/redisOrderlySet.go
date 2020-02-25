package mRedis

import (
	"fmt"
	"markV5/mError"
	"markV5/mFace"
)

/*
经典使用场景:
1.排行榜系统:例如视频点赞数目排行，视频主需要通过有序集合来确保去重，通过 zIncrBy 来增加点赞数，通过 zRem 来
取消作弊的视频主，通过 zRevRangeByRank 来进行排行
*/
/*
有序集合特点:
1.元素不可重复，与集合类似
2.有序，但不支持下标查询，支持权重查询，通过给元素分配score来排序，score可重复
*/

type ZRangeByScore_Limit string
type ZInterStore_Aggregate string

const (
	ZRangeByScore_In            = ""
	ZRangeByScore_NotIn         = "("
	ZRangeByScore_Infinitesimal = "-inf"
	ZRangeByScore_Infinity      = "+inf"
)

const (
	ZInterStore_Sum ZInterStore_Aggregate = "sum"
	ZInterStore_Max ZInterStore_Aggregate = "max"
	ZInterStore_Min ZInterStore_Aggregate = "min"
)

// 向 key 中添加 members ， members 由权重分数和值组成
func (r *MRedis) ZAdd(key string, members map[int]interface{}) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for score, member := range members {
		params = append(params, score)
		params = append(params, member)
	}
	return r.Exec("ZADD", key, params...)
}

// 统计 key 的元素个数
func (r *MRedis) ZCard(key string) (interface{}, mFace.MError) {
	return r.Exec("ZCARD", key)
}

// 统计 key 的元素 member 的权重分数
func (r *MRedis) ZScore(key string, member interface{}) (interface{}, mFace.MError) {
	return r.Exec("ZSCORE", key, member)
}

// 统计 key 的元素 member 的排名，从低到高返回排名，排名从0开始计算
func (r *MRedis) ZRank(key string, member interface{}) (interface{}, mFace.MError) {
	return r.Exec("ZRANK", key, member)
}

// 统计 key 的元素 member 的排名，从高到低返回排名，排名从0开始计算
func (r *MRedis) ZRevRank(key string, member interface{}) (interface{}, mFace.MError) {
	return r.Exec("ZREVRANK", key, member)
}

// 删除 key 的 members元素
func (r *MRedis) ZRem(key string, members ...interface{}) (interface{}, mFace.MError) {
	return r.Exec("ZREM", key, members...)
}

// 增加 key 的 member 元素的权重分数 increment
func (r *MRedis) ZIncrBy(key string, increment int, member interface{}) (interface{}, mFace.MError) {
	return r.Exec("ZINCRBY", key, increment, member)
}

// 返回 key 指定排名范围 start - end 的元素，从低到高，包含end元素
func (r *MRedis) ZRange(key string, start, end int, withScores bool) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	params = append(params, start)
	params = append(params, end)
	if withScores {
		params = append(params, "withscores")
	}
	return r.Exec("ZRANGE", key, params...)
}

// 返回 key 指定排名范围 start - end 的元素，从高到低，包含end元素
func (r *MRedis) ZRevRange(key string, start, end int, withScores bool) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	params = append(params, start)
	params = append(params, end)
	if withScores {
		params = append(params, "withscores")
	}
	return r.Exec("ZRevRANGE", key, params...)
}

// 返回 key 中指定权重分数范围 min - max 的元素，排名从低到高
func (r *MRedis) ZRangeByScore(key string,
	minLimit ZRangeByScore_Limit, min int,
	maxLimit ZRangeByScore_Limit, max int, withScores,
	limit bool, offset, count int) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	switch minLimit {
	case ZRangeByScore_NotIn:
		params = append(params, fmt.Sprintf("(%d", min))
	case ZRangeByScore_Infinitesimal:
		params = append(params, "-inf")
	default:
		params = append(params, min)
	}
	switch maxLimit {
	case ZRangeByScore_NotIn:
		params = append(params, fmt.Sprintf("(%d", max))
	case ZRangeByScore_Infinity:
		params = append(params, "+inf")
	default:
		params = append(params, max)
	}
	if withScores {
		params = append(params, "withscores")
	}
	if limit {
		params = append(params, "limit")
		params = append(params, offset)
		params = append(params, count)
	}
	return r.Exec("ZRANGEBYSCORE", key, params...)
}

// 返回 key 中指定权重分数范围 min - max 的元素，排名从高到低
func (r *MRedis) ZRevRangeByScore(key string,
	maxLimit ZRangeByScore_Limit, max int,
	minLimit ZRangeByScore_Limit, min int, withScores,
	limit bool, offset, count int) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	switch maxLimit {
	case ZRangeByScore_NotIn:
		params = append(params, fmt.Sprintf("(%d", max))
	case ZRangeByScore_Infinity:
		params = append(params, "+inf")
	default:
		params = append(params, max)
	}
	switch minLimit {
	case ZRangeByScore_NotIn:
		params = append(params, fmt.Sprintf("(%d", min))
	case ZRangeByScore_Infinitesimal:
		params = append(params, "-inf")
	default:
		params = append(params, min)
	}
	if withScores {
		params = append(params, "withscores")
	}
	if limit {
		params = append(params, "limit")
		params = append(params, offset)
		params = append(params, count)
	}
	return r.Exec("ZREVRANGEBYSCORE", key, params...)
}

// 返回 key 指定权重分数范围 min - max 内的个数
func (r *MRedis) ZCount(key string, min, max int) (interface{}, mFace.MError) {
	return r.Exec("ZCOUNT", key, min, max)
}

// 删除 key 指定排名范围 start - end 的元素，包含end
func (r *MRedis) ZRemRangeByRank(key string, start, end int) (interface{}, mFace.MError) {
	return r.Exec("ZREMRANGEBYRANK", key, start, end)
}

// 删除 key 指定权重分数范围 min - max 的元素，包含end
func (r *MRedis) ZRemRangeByScore(key string, min, max int) (interface{}, mFace.MError) {
	return r.Exec("ZREMRANGEBYScore", key, min, max)
}

// ------------------------ 集合间操作

// 根据 keys 中各个key及其权重，按照 aggregate 指定的方式交集计算到指定的 desKey 中
func (r *MRedis) ZInterStore(desKey string, keys map[string]float64, aggregate ZInterStore_Aggregate) (interface{}, mFace.MError) {
	if len(keys) == 0 {
		return nil, mError.ErrorParam
	}
	params := make([]interface{}, 0)
	params = append(params, len(keys))
	weights := make([]interface{}, 0)
	for key, weight := range keys {
		params = append(params, key)
		weights = append(weights, weight)
	}
	params = append(params, "weights")
	params = append(params, weights...)
	params = append(params, "aggregate")
	switch aggregate {
	case ZInterStore_Sum:
		params = append(params, "sum")
	case ZInterStore_Max:
		params = append(params, "max")
	case ZInterStore_Min:
		params = append(params, "min")
	default:
		params = params[:len(params)-1]
	}
	return r.Exec("ZINTERSTORE", desKey, params...)
}

// 根据 keys 中各个key及其权重，按照 aggregate 指定的方式并集计算到指定的 desKey 中
func (r *MRedis) ZUnionStore(desKey string, keys map[string]float64, aggregate ZInterStore_Aggregate) (interface{}, mFace.MError) {
	if len(keys) == 0 {
		return nil, mError.ErrorParam
	}
	params := make([]interface{}, 0)
	params = append(params, len(keys))
	weights := make([]interface{}, 0)
	for key, weight := range keys {
		params = append(params, key)
		weights = append(weights, weight)
	}
	params = append(params, "weights")
	params = append(params, weights...)
	params = append(params, "aggregate")
	switch aggregate {
	case ZInterStore_Sum:
		params = append(params, "sum")
	case ZInterStore_Max:
		params = append(params, "max")
	case ZInterStore_Min:
		params = append(params, "min")
	default:
		params = params[:len(params)-1]
	}
	return r.Exec("ZUNIONSTORE", desKey, params...)
}
