package connect

type ConnecError struct {
	Err     error
	ErrCode int
}

func MakeConnectErr(err error, errCode int) ConnecError {
	return ConnecError{
		Err:     err,
		ErrCode: errCode,
	}
}

type Connector interface {
	SetConfig(any)
	Request() ([]byte, ConnecError)
	ConnectorID() string
}
