package utils

import "strings"

type Mapper map[string]interface{}

func ReplaceQuote(key string) string {
	return strings.ReplaceAll(key, `"`, "")
}

func (m Mapper) Find(key string) (interface{}, bool) {
	keys := strings.Split(ReplaceQuote(key), ".")
	return NestedMapLookup(m, keys...)
}

func (m Mapper) FindWithDefault(key string, default_ interface{}) interface{} {
	if v, ok := m.Find(key); ok {
		return v
	}
	if v, ok := default_.(string); ok {
		return ReplaceQuote(v)
	}
	return default_
}

func NestedMapLookup(m map[string]interface{}, ks ...string) (interface{}, bool) {
	var ok bool
	var val interface{}

	// 入力の検証
	if len(ks) == 0 {
		return nil, false
	}
	if ks[0] == "" {
		return nil, false
	}

	// 最初のキーで値を取得する。
	val, ok = m[ks[0]]
	if !ok {
		return nil, false
	}

	// 最後のキーの場合、値を返す。
	if len(ks) == 1 {
		return val, true
	}

	// 値がマップの場合、再帰的にネストされたマップを探索する。
	if m, ok := val.(map[string]interface{}); ok {
		return NestedMapLookup(m, ks[1:]...)
	}

	// 値がマップでない場合、エラーを返す。
	return nil, false
}
