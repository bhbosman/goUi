package cmSlide

type DisconnectConnection struct {
	ConnectionId   string
	ConnectionName string
}

func NewDisconnectConnection(connectionId string, connectionName string) *DisconnectConnection {
	return &DisconnectConnection{ConnectionId: connectionId, ConnectionName: connectionName}
}
