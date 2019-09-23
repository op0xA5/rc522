package rc522

type Err int

const (
	noErr = Err(iota)
	timeoutErr
	collisionErr
	badCrcErr
)

func (err Err) Error() string {
	switch err {
	case timeoutErr:
		return "timeout"
	case collisionErr:
		return "collision"
	case badCrcErr:
		return "bad crc"
	}
	return ""
}

func IsTimeout(err error) bool {
	if e, ok := err.(Err); ok {
		return e == timeoutErr
	}
	return false
}

func IsCollision(err error) bool {
	if e, ok := err.(Err); ok {
		return e == collisionErr
	}
	return false
}
