package util

// 雪花算法生成id
import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

// WorkerID worker
const WorkerID = 10

var (
	sf *sonyflake.Sonyflake
)

func init() {
	st := sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return WorkerID, nil
		},
	}
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("not created")
	}
}

// GetNextID 获取雪花算法生成的id
func GetNextID() (uint64, error) {
	return sf.NextID()
}

// MustGetNextID 获取雪花算法生成的id 有错直接panic
func MustGetNextID() uint64 {
	res, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return res
}

// NewUUID 创建新的uuid
func NewUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
