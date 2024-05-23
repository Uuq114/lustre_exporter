package collector

import (
	"math/rand"
	"time"
)

type ExampleMetrics struct {
	randomFloat  float32
	randomString string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func exportMetric() map[string]any {
	var exampleMetric ExampleMetrics
	rand.Seed(time.Now().UnixMilli())
	exampleMetric.randomFloat = rand.Float32()
	exampleMetric.randomString = getRandomString(rand.Intn(10))

	return map[string]any{
		"RandomFloat":  exampleMetric.randomFloat,
		"RandomString": exampleMetric.randomString,
	}
}
