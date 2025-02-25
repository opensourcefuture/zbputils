package pool

import (
	"errors"
	"time"

	"github.com/FloatTech/zbputils/process"
	"github.com/fumiama/go-registry"
)

type item struct {
	name string
	u    string
}

// newItem 唯一标识文件名 文件链接
func newItem(name, u string) (*item, error) {
	if len(name) > 126 {
		return nil, errors.New("name too long")
	}
	if len(u) > 126 {
		return nil, errors.New("url too long")
	}
	return &item{name: name, u: u}, nil
}

// getItem 唯一标识文件名
func getItem(name string) (*item, error) {
	reg := registry.NewRegReader("reilia.westeurope.cloudapp.azure.com:35354", "fumiama")
	err := reg.ConnectIn(time.Second * 4)
	if err != nil {
		return nil, err
	}
	u, err := reg.Get(name)
	defer reg.Close()
	if err != nil {
		return nil, err
	}
	return &item{name: name, u: u}, nil
}

// push 推送 item
func (t *item) push(key string) (err error) {
	r := registry.NewRegedit("reilia.westeurope.cloudapp.azure.com:35354", "fumiama", key)
	err = r.ConnectIn(time.Second * 4)
	if err != nil {
		return
	}
	for i := 0; i < 8; i++ {
		err = r.Set(t.name, t.u)
		if err == nil {
			break
		}
		process.SleepAbout1sTo2s() // 随机退避
	}
	_ = r.Close()
	return
}
