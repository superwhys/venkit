package shared

var (
	ConsulAddr func() string
)

func GetConsulAddress() string {
	return ConsulAddr()
}
