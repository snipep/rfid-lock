package mqtt

const (
	Broker = "broker.emqx.io:1883"
)

// List of topics to subscribe to
var Topics = []string{
	"rfid/auth",
	"rfid/auth/status",
	"register/user/valid",
	"register/user/invalid",
}
