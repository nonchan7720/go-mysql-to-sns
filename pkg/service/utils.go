package service

import "github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"

func findTarget[T any](mp map[string]T, payload interfaces.SendPayload) (T, bool) {
	if v, ok := mp[payload.GetTable()]; ok {
		return v, ok
	} else {
		var v T
		return v, false
	}
}

func FindTarget[T any](mp map[string]T, payload interfaces.SendPayload) T {
	if v, ok := findTarget(mp, payload); ok {
		return v
	}
	panic("Target not found.")
}

func IsTarget[T any](mp map[string]T, payload interfaces.SendPayload) bool {
	_, ok := findTarget(mp, payload)
	return ok
}
