package main

import (
	"fmt"
	"log"

	"github.com/neuralmagic/distributed-kv-cache/pkg/kvcacheindexer"
)

func main() {
	indexer := kvcacheindexer.NewKVCacheIndexer()
	prompt := []string{"The", "sky", "is"}
	model := kvcacheindexer.ModelInfo{Name: "llm", Version: "v1"}
	allowedPods := []string{"pod-a"} // Filter list for relevance

	fmt.Println("Using strategy:", indexer.Scorer.Strategy())

	// Step 1: Run pipeline before prefix update
	pods, err := indexer.RunPrefixAwarePipeline(prompt, model, allowedPods)
	if err != nil {
		log.Fatalf("Initial RunPrefixAwarePipeline failed: %v", err)
	}
	fmt.Println("Before Update: Pods matched =", pods)

	// Step 2: Update prefix cache to map prompt to pod-a
	hashes := indexer.PrefixCache.GetPrefixHashes(prompt)
	for _, h := range hashes {
		indexer.PrefixCache.UpdatePodPrefix(h, "pod-a")
	}

	// Step 3: Run pipeline again — should now return pod-a
	updatedPods, err := indexer.RunPrefixAwarePipeline(prompt, model, allowedPods)
	if err != nil {
		log.Fatalf("Post-update RunPrefixAwarePipeline failed: %v", err)
	}

	fmt.Println("After Update: Pods matched =")
	for _, pod := range updatedPods {
		fmt.Printf(" - %s (Score: %.2f)\n", pod.Name, pod.Score)
	}
}
