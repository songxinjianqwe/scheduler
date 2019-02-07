package client
import (
	"github.com/songxinjianqwe/go-practice/concurrent/timer/common"
)
type Client struct {

}
/**
构造一个Client实例
 */
func  NewClient() (c *Client) {
	return nil
}

func (this *Client) list() []common.Task {
	return nil
}

func (this *Client) submit(task *common.Task) error {
	return nil
}

func (this *Client) watch(id string) error {
	return nil
}


