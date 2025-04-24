package j60

import (
	"errors"
	"github.com/brutella/can"
	"github.com/injoyai/base/maps/wait/v2"
	"github.com/injoyai/conv"
	"github.com/injoyai/logs"
	"math"
	"time"
)

type Frame = can.Frame
type any = interface{}

const (
	DefaultInterval = time.Millisecond * 100
	DefaultAccuracy = 5
)

func Dial(name string) (*Client, error) {
	c, err := can.NewBusForInterfaceWithName(name)
	if err != nil {
		return nil, err
	}
	wait := wait.New(time.Second * 2)
	c.SubscribeFunc(func(f can.Frame) {
		var result any
		id, cmd, ok := DID(f.ID)
		key := conv.String(id) + conv.String(cmd)
		if !ok {
			result = errors.New("失败")
		} else if cmd == CONTROL_MOTOR {
			res := DecodeControlResp(f.Data)
			result = res
			logs.Readf("id: %d, 角度: %.2f, 扭矩: %.2f, 温度: %.2f\n", id, res.Position, res.Velocity, res.Temperature)
		}
		wait.Done(key, result)
	})
	go c.ConnectAndPublish()
	return &Client{
		Bus:  c,
		wait: wait,
	}, nil
}

type Client struct {
	*can.Bus
	wait *wait.Entity
}

// Enable 使能
func (this *Client) Enable(id uint8) error {
	return this.WriteCommand(id, ENABLE_MOTOR)
}

// Disable 失能
func (this *Client) Disable(id uint8) error {
	return this.WriteCommand(id, DISABLE_MOTOR)
}

// Reset 重置
func (this *Client) Reset(id uint8) error {
	return this.WriteCommand(id, RESET_MOTOR)
}

// Calibrate 校准
func (this *Client) Calibrate(id uint8) error {
	return this.WriteCommand(id, CALIBRATE_START)
}

// SetZero 设置零点
func (this *Client) SetZero(id uint8) error {
	return this.WriteCommand(id, SET_HOME)
}

func (this *Client) WriteCommand(id, cmd uint8) error {
	logs.Writef("id: %d, cmd: %d, can_id: %d\n", id, cmd, ID(id, cmd))
	err := this.Bus.Publish(can.Frame{ID: ID(id, cmd)})
	if err != nil {
		return err
	}
	_, err = this.wait.Wait(conv.String(id) + conv.String(cmd))
	return err
}

// Control 控制电机
func (this *Client) Control(id uint8, req *ControlReq) (*ControlResp, error) {
	bs, err := req.Bytes()
	if err != nil {
		return nil, err
	}
	logs.Writef("id: %d, 角度: %.2f, 刚度: %.2f, 阻尼: %.2f\n", id, req.Position, req.Kp, req.Kd)
	f := can.Frame{
		ID:     ID(id, CONTROL_MOTOR),
		Length: uint8(len(bs)),
		Data:   *bs,
	}
	err = this.Bus.Publish(f)
	if err != nil {
		return nil, err
	}
	val, err := this.wait.Wait(conv.String(id) + conv.String(CONTROL_MOTOR))
	if err != nil {
		return nil, err
	}
	return val.(*ControlResp), nil
}

// ControlPosition 控制角度
func (this *Client) ControlPosition(id uint8, p, kp, kd float64) (*ControlResp, error) {
	return this.Control(id, &ControlReq{
		Position: p,
		Kp:       kp,
		Kd:       kd,
	})
}

func (this *Client) ControlPositionWait(id uint8, p, kp, kd float64) error {
	for {
		resp, err := this.ControlPosition(id, p, kp, kd)
		if err != nil {
			return err
		}
		if math.Abs(p-resp.Position) <= DefaultAccuracy {
			return nil
		}
		<-time.After(DefaultInterval)
	}
}

func (this *Client) WriteFrame(f can.Frame) error {
	return this.Bus.Publish(f)
}

func (this *Client) Close() error {
	return this.Bus.Disconnect()
}
