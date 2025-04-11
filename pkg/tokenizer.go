package pkg

import (
	"fmt"
	"github.com/daulet/tokenizers"
)

func Tokenize(input, modelName, hfToken string) ([]uint32, error) {
	tk, err := tokenizers.FromPretrained(modelName, tokenizers.WithAuthToken(hfToken))
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
