package zinf

type ZinfRequest interface {
	GetConnection() ZinfConnection
	GetData() []byte
}
