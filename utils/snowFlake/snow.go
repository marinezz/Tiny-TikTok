package snowFlake

import (
	"errors"
	"sync"
	"time"
)

const (
	// 当前时间 2023-08-07 10:47:09  毫秒级别
	epoch int64 = 1691376429015
	// 序列号位数
	numberBit int8 = 12
	// 服务ID位数
	serveIdBit int8 = 5
	// 机器ID位数
	machineIdBit int8 = 5
	// 服务ID的偏移量
	serveIdShift int8 = numberBit
	// 机器ID的偏移量
	machineIdShift int8 = numberBit + serveIdBit
	// 时间戳的偏移量
	timestampShift int8 = numberBit + serveIdBit + machineIdBit
	// 服务ID的最大值 31
	serverIdMax int64 = -1 ^ (-1 << serveIdBit)
	// 机器ID的最大值 31
	machineIdMax int64 = -1 ^ (-1 << machineIdBit)
	// 序列号的最大值 4095
	numberMax int64 = -1 ^ (-1 << numberBit)
)

type SnowFlake struct {
	// 每次生产一个id都是一个原子操作
	lock sync.Mutex
	// 时间戳、机器ID、服务ID、序列号
	timestamp int64
	machineId int64
	serveId   int64
	number    int64
}

// NewSnowFlake 构造函数，传入机器ID和服务ID
func NewSnowFlake(machineId int64, serveId int64) (*SnowFlake, error) {
	if machineId < 0 || machineId > machineIdMax {
		return nil, errors.New("mechineId超出限制")
	}
	if serveId < 0 || serveId > serverIdMax {
		return nil, errors.New("serveId超出限制")
	}
	return &SnowFlake{
		timestamp: 0,
		machineId: machineId,
		serveId:   serveId,
		number:    0,
	}, nil
}

func (snow *SnowFlake) NextId() int64 {
	// 原子操作
	snow.lock.Lock()
	defer snow.lock.Unlock()

	// 获取当前时间戳
	now := time.Now().UnixMilli()

	// 如果时间戳还是当前时间错，则序列号增加
	if now == snow.timestamp {
		snow.number++
		// 如果超过了序列号的最大值，则更新时间戳
		if snow.number > numberMax {
			for now <= snow.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		snow.number = 0
		snow.timestamp = now
	}

	// 拼接最后的结果，将不同不服的数值移到指定位置
	id := (snow.timestamp-epoch)<<timestampShift | (snow.machineId << machineIdShift) |
		(snow.serveId << serveIdShift) | snow.number

	return id
}
