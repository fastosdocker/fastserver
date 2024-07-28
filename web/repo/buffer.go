package repo

import "sync"

var (
	buffs      bufferList
	buffsMutex sync.RWMutex
)

func init() {
	buffs = make(bufferList)
}

const maxRecord = 3600 // 每个容器最大存放条数

// 存放docker资源利用率的环形队列
type Buffer struct {
	CID   string
	data  [maxRecord][]string // 使用固定大小数组代替切片
	head  int                 // 记录队列头的位置
	tail  int                 // 记录队列尾的位置
	count int                 // 队列中元素的数量
	mu    sync.RWMutex
}

// 初始化buffer，设置队列头和尾的初始位置
func newBuffer(cid string) *Buffer {
	buffsMutex.Lock()
	defer buffsMutex.Unlock()
	if buffs == nil {
		buffs = make(bufferList)
	}
	addBuffer(cid)
	return &Buffer{
		CID:   cid,
		head:  0,
		tail:  0,
		count: 0,
	}
}

// 向buffer中追加新记录
func (b *Buffer) append(newRecord []string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 如果队列已满，覆盖最旧的记录
	if b.count == maxRecord {
		b.head = (b.head + 1) % maxRecord
		b.count--
	} else if b.tail == maxRecord {
		b.tail = 0
	}

	// 添加新记录
	b.data[b.tail] = newRecord
	b.tail++
	b.count++
}

// 读取最新的记录
func (b *Buffer) read() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.count > 0 {
		return b.data[(b.tail-1+maxRecord)%maxRecord]
	}
	return nil
}

// bufferList用于管理多个buffers
type bufferList map[string]*Buffer

// 根据CID查找buffer
func GetBuffer(id string) (*Buffer, bool) {
	buffsMutex.RLock()
	defer buffsMutex.RUnlock()
	b, ok := buffs[id]
	return b, ok
}

// 添加新的buffer到bufferList
func addBuffer(cid string) {
	buffsMutex.Lock()
	defer buffsMutex.Unlock()
	buffs[cid] = newBuffer(cid)
}
