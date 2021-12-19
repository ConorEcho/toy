package toy

import "fmt"

type hashStorage struct {
	hashMap map[string]HandlerChain
}

func NewHashStorage() routeStorage {
	return &hashStorage{make(map[string]HandlerChain)}
}

func (t *hashStorage) Store(method string, route string, handlers HandlerChain) {
	t.hashMap[fmt.Sprintf("%s-%s", method, route)] = handlers

}

func (t *hashStorage) GetHandlers(method string, route string) HandlerChain {
	return t.hashMap[fmt.Sprintf("%s-%s", method, route)]
}
