package kvcacheindexer

// KVScoringStrategy defines the strategy used to score pods for KV cache block reuse.
type KVScoringStrategy string

const (
	// LongestPrefixMatch Score by longest consecutive match from start.
	LongestPrefixMatch KVScoringStrategy = "LongestPrefix"
	// HighestBlockHit Score by highest block index hit.
	HighestBlockHit KVScoringStrategy = "HighestBlockHit"
	// CoverageBasedMatching Score by total number of blocks hit.
	CoverageBasedMatching KVScoringStrategy = "CoverageBased"
	//  UserDefinedStrategy Score via user-defined strategy.
	UserDefinedStrategy KVScoringStrategy = "UserDefined"
)

// KVScorer defines the interface for implementing a KV block scoring strategy.
type KVScorer interface {
	Strategy() KVScoringStrategy
	Score(blockKeys []string, hitmap map[string]string) ([]Pod, error)
}

// ------------------------
// LongestPrefixScorer
// ------------------------

// LongestPrefixScorer scores based on longest consecutive block matches from the beginning.
type LongestPrefixScorer struct{}

// Strategy returns the strategy type: LongestPrefixMatch.
func (s *LongestPrefixScorer) Strategy() KVScoringStrategy {
	return LongestPrefixMatch
}

// Score implements the longest prefix scoring logic.
func (s *LongestPrefixScorer) Score(blockKeys []string, hitmap map[string]string) ([]Pod, error) {
	longestIndex := make(map[string]int)

	if len(blockKeys) == 0 {
		return nil, nil
	}

	prevPod, ok := hitmap[blockKeys[0]]
	if !ok {
		return nil, nil
	}
	count := 1

	for i := 1; i < len(blockKeys); i++ {
		pod, ok := hitmap[blockKeys[i]]
		if !ok {
			break
		}
		if pod == prevPod {
			count++
		} else {
			if count > longestIndex[prevPod] {
				longestIndex[prevPod] = count
			}
			prevPod = pod
			count = 1
		}
	}

	if count > longestIndex[prevPod] {
		longestIndex[prevPod] = count
	}

	return convertScoresToPods(longestIndex), nil
}

// ------------------------
// HighestBlockHitScorer
// ------------------------

// HighestBlockHitScorer scores based on the highest-indexed block hit for each pod.
type HighestBlockHitScorer struct{}

// Strategy returns the strategy type: HighestBlockHit.
func (s *HighestBlockHitScorer) Strategy() KVScoringStrategy {
	return HighestBlockHit
}

// Score implements the highest block hit scoring logic.
func (s *HighestBlockHitScorer) Score(blockKeys []string, hitmap map[string]string) ([]Pod, error) {
	maxIndex := make(map[string]int)

	for idx, k := range blockKeys {
		pod, ok := hitmap[k]
		if !ok {
			continue
		}
		maxIndex[pod] = idx
	}

	return convertScoresToPods(maxIndex), nil
}

// ------------------------
// CoverageBasedScorer
// ------------------------

// CoverageBasedScorer scores based on total number of blocks hit (coverage).
type CoverageBasedScorer struct{}

// Strategy returns the strategy type: CoverageBasedMatching.
func (s *CoverageBasedScorer) Strategy() KVScoringStrategy {
	return CoverageBasedMatching
}

// Score implements the coverage-based scoring logic.
func (s *CoverageBasedScorer) Score(blockKeys []string, hitmap map[string]string) ([]Pod, error) {
	coverage := make(map[string]int)

	for _, k := range blockKeys {
		pod, ok := hitmap[k]
		if !ok {
			continue
		}
		coverage[pod]++
	}

	return convertScoresToPods(coverage), nil
}

// ------------------------
// Shared helper
// ------------------------

// convertScoresToPods converts a map of pod name to score into a slice of Pod structs.
func convertScoresToPods(scoreMap map[string]int) []Pod {
	scored := make([]Pod, 0, len(scoreMap))
	for pod, score := range scoreMap {
		scored = append(scored, Pod{
			Name:  pod,
			Score: float64(score),
		})
	}
	return scored
}
