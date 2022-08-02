all: IConnectionReactorFactory


IConnectionReactorFactory:
	mockgen -package netListener -generateWhat mockgen -destination ConnMock.go net Conn,Addr,Listener
	mockgen -package netListener -generateWhat mockgen -destination IListenerAcceptMock.go . IListenerAccept







