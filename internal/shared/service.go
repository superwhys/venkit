package shared

var (
	ServiceName func() string
)

func GetServiceName() string {
	return ServiceName()
}
