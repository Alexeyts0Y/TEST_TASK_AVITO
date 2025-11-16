package utils

import (
	"math/rand"
	"time"
)

func ChooseRandomCandidates(candidates []string, count int) []string {
	if len(candidates) <= count {
		return candidates
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := random.Perm(len(candidates))
	chosen := make([]string, 0, count)
	for i := 0; i < count; i++ {
		chosen = append(chosen, candidates[perm[i]])
	}
	return chosen
}
