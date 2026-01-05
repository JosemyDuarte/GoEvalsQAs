package eval

import (
	"fmt"
	"strings"
)

// FlattenContext turns the complex HotpotQA context array into a readable string.
func FlattenContext(ctx [][]interface{}) string {
	var builder strings.Builder

	for _, doc := range ctx {
		if len(doc) < 2 {
			continue
		}

		title, ok := doc[0].(string)
		if !ok {
			continue
		}

		sentencesRaw, ok := doc[1].([]interface{})
		if !ok {
			continue
		}

		builder.WriteString(fmt.Sprintf("Document [%s]: ", title))

		for _, s := range sentencesRaw {
			if str, ok := s.(string); ok {
				builder.WriteString(str)
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
