package vecstore

import (
	"lazyollama/model"
	"math"
	"sort"
)

func cosineSim(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func FindTopKRelevant(
	docs []model.Embedding,
	queryEmbedding []float64,
	topK int,
) []model.Embedding {
	sort.Slice(docs, func(i, j int) bool {
		return cosineSim(
			queryEmbedding,
			docs[i].Emb,
		) > cosineSim(
			queryEmbedding,
			docs[j].Emb,
		)
	})
	if topK > len(docs) {
		topK = len(docs)
	}
	return docs[:topK]
}
