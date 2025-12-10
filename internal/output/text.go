package output

import (
	"fmt"
	"io"
)

// WriteText はURLリストをテキスト形式で出力する
func WriteText(w io.Writer, urls []string) error {
	for _, url := range urls {
		if _, err := fmt.Fprintln(w, url); err != nil {
			return err
		}
	}
	return nil
}
