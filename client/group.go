package client

type GroupSyncReader interface {
	Read(chan *GroupSyncReadResult)
}

type GroupSyncReadResult struct {
	Error   error
	Results []*SyncReadResult
}

type SyncReadResult struct {
	Error  error
	ID     string
	Result []byte
}
