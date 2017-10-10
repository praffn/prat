package prat

const (
	DefaultPort = 9876
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
