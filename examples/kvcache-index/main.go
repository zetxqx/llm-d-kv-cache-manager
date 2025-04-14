package main

import (
	"time"

	"github.com/neuralmagic/distributed-kv-cache/pkg/kvcache"

	"golang.org/x/net/context"
	"k8s.io/klog/v2"
)

/*
Refer to docs/phase1-setup.md

In Redis:
1) "vllm@<modelName>@1@0@968e5b25ee324e1207bb88f4b4ad208731cbff985ac3ec6ce452b70aa54df4b4"
2) "vllm@<modelName>@1@0@a9d880fb27ad4f27cd8affaf3607796ccb4f09bd2658b736766d0452c612cc93"
*/

//nolint:lll // need prompt as-is, chunking to string concatenation is too much of a hassle
const prompt = `lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst.`
const modelName = "mistralai/Mistral-7B-Instruct-v0.2"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := klog.FromContext(ctx)

	// TODO: create a configuration module with default config option
	kvCacheIndexer, err := kvcache.NewKVCacheIndexer(kvcache.NewDefaultConfig())
	if err != nil {
		logger.Error(err, "failed to init Indexer")
	}

	logger.Info("created Indexer")

	go kvCacheIndexer.Run(ctx)
	logger.Info("started Indexer")

	// Get pods for the prompt
	pods, err := kvCacheIndexer.GetPodScores(ctx, prompt, modelName, nil)
	if err != nil {
		logger.Error(err, "failed to get pod scores")
		return
	}

	// Print the pods - should be empty because no tokenization
	logger.Info("got pods", "pods", pods)

	// Sleep 3 secs
	time.Sleep(3 * time.Second)

	// Get pods for the prompt
	pods, err = kvCacheIndexer.GetPodScores(ctx, prompt, modelName, nil)
	if err != nil {
		logger.Error(err, "failed to get pod scores")
		return
	}

	// Print the pods - should be empty because no tokenization
	logger.Info("got pods", "pods", pods)

	// Cancel the context
	cancel()
}
