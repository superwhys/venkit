package shared

var (
	UseConsul  func() bool
	ConsulAddr func() string
)

func GetIsUseConsul() bool {
	if UseConsul == nil {
		return false
	}

	return UseConsul()
}

func GetConsulAddress() string {
	if ConsulAddr == nil {
		return ""
	}
	return ConsulAddr()
}
