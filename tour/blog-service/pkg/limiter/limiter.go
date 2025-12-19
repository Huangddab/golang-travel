package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 定义限流器所必须要的方法
type LimiterIface interface {
	Key(c *gin.Context) string                           // 获取对应的限流器的键值对
	GetBucket(key string) (*ratelimit.Bucket, bool)      // 获取令牌桶
	AddBuckets(rules ...LimiterBucketsRule) LimiterIface // 新增多个令牌桶
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

type LimiterBucketsRule struct {
	Key          string        // 自定义键值对名称
	FillInterval time.Duration // 间隔多久时间放N个令牌
	Capacity     int64         // 令牌桶的容量
	Quantum      int64         // 每次到达间隔时间过放的具体令牌数量
}
