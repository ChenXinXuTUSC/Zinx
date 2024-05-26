package zinf

type ZinfConnManager interface {
	Add(conn ZinfConnection)
	Remove(conn ZinfConnection)
	Get(connId uint32) (ZinfConnection, error)
	GetConnNum() int
	ClearAll()
}
