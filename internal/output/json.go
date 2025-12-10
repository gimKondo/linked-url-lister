package output

import (
	"encoding/json"
	"io"
)

// JSONOutput はJSON出力の構造体
type JSONOutput struct {
	URLs   []string    `json:"urls"`
	Stats  JSONStats   `json:"stats"`
	Errors []JSONError `json:"errors"`
}

// JSONStats は統計情報
type JSONStats struct {
	TotalVisited    int     `json:"total_visited"`
	TotalFiltered   int     `json:"total_filtered"`
	DurationSeconds float64 `json:"duration_seconds"`
}

// JSONError はエラー情報
type JSONError struct {
	URL        string `json:"url"`
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

// WriteJSONOutput はJSONOutputをJSON形式で出力する
func WriteJSONOutput(w io.Writer, output *JSONOutput) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
