package ads

import "errors"

type Result uint32

func (this Result) Err() error {
	if this != 0 {
		return errors.New("失败")
	}
	return nil
}
