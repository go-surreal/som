//go:build embed

package internal

type QueryResult[M any] struct {
	Result []M    `json:"result"`
	Status string `json:"status"`
	Time   string `json:"time"`
}
