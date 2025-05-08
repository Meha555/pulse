package iface

type IMap[K comparable, V any] interface {
	Store(key K, value any)
	Load(key K) (value V, exists bool)
}
