package voice

import "errors"

var (
	GatewayMissingValuesError = errors.New("some of the required values required to start connection are missing")
	ConnectionClosedError     = errors.New("connection is closed")
	UDPHolepunchFailed        = errors.New("failed to perform UDP holepunch in 10 tries")
	NotConnected              = errors.New("session is not yet connected, or is reconnecting")
)
