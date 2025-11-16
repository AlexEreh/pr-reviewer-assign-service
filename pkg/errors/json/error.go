package json

const (
	keyCode       = "Code"
	keyMessage    = "Message"
	keyCause      = "Cause"
	keyStackTrace = "StackTrace"
)

type Error struct {
	Error       error
	Marshaler   []MarshalerOption
	Unmarshaler []UnmarshalerOption
}
