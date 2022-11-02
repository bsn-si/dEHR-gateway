package publisher

type Subscriber interface {
	Notify(interface{})
	Disable()
	Name() string
}
