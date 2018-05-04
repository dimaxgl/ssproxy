package memory

import "github.com/dimaxgl/ssproxy/store"

type memoryStory struct {

}

func (memoryStory) Valid(user, password string) bool {
	panic("implement me")
}

func (memoryStory) Add(user, password string) error {
	panic("implement me")
}

func (memoryStory) Initialize(params map[string]interface{}) (store.Store, error) {
	panic("implement me")
}
