package limiter

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() LimiterIface {
	return &MethodLimiter{ // 返回指针
		Limiter: &Limiter{limiterBuckets: make(map[string]*ratelimit.Bucket)},
	}
}

func (l *MethodLimiter) Key(c *gin.Context) string { // 接收者改为指针
	uri := c.Request.RequestURI
	index := strings.Index(uri, "?")
	if index == -1 {
		return uri
	}

	return uri[:index]
}

func (l *MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) { // 接收者改为指针
	bucket, ok := l.limiterBuckets[key]
	return bucket, ok
}

func (l *MethodLimiter) AddBuckets(rules ...LimiterBucketsRule) LimiterIface { // 接收者改为指针
	for _, rule := range rules {
		if _, ok := l.limiterBuckets[rule.Key]; !ok {
			l.limiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}

	return l
}
