package utils

import "net/url"

func PtrX[T any](v T) *T {
	return &v
}

func SoftPathUnescape(raw string) string {
	res, err := url.PathUnescape(raw) // 非 ASCII 的字符要做额外处理
	if err != nil {
		return raw // 当出错时就返回原始的
	}
	return res
}