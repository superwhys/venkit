package shared

var (
	ServiceName func() string
)

func GetServiceName() string {
	if ServiceName == nil {
		return ""
	}
	return ServiceName()
}
