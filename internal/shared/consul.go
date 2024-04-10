package shared

var (
	ConsulAddr func() string
)

func GetConsulAddress() string {
	if ConsulAddr == nil {
		return ""
	}
	return ConsulAddr()
}
