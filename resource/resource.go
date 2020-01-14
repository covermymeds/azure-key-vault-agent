package resource

type Resource interface {
	Map() map[string]interface{}
	String() string
}
