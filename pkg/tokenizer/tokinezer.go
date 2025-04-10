package tokenizer

import (
	"fmt"
	"os"

	"github.com/daulet/tokenizers"

	"github.com/sirupsen/logrus"
)

const DefaultEncoding = "cl100k_base"

type Tokenizer struct {
	cfg    tokenizers.TokenizerConfigOption
	logger *logrus.Entry
}

// NewTokenizer creates a Tiktoken tokenizer.
func NewTokenizer() *Tokenizer {
	logger := logrus.WithField("component", "tokenizer")

	cfg := tokenizers.WithAuthToken(os.Getenv("HF_TOKEN")) // Todo- use cache dir
	return &Tokenizer{
		logger: logger,
		cfg:    cfg,
	}
}

// Encode converts a string into token IDs.
func (t *Tokenizer) Encode(input, modelName string) ([]uint32, error) {
	tk, err := tokenizers.FromPretrained(modelName, t.cfg)
	if err != nil {
		return nil, err
	}

	defer func(tk *tokenizers.Tokenizer) {
		err := tk.Close()
		if err != nil {
			fmt.Println("Error closing tokenizer:", err)
		}
	}(tk)

	ids, _ := tk.Encode(input, true)
	return ids, nil
}
