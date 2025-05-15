package rabbitmq

type PublishOption struct {
	Exchange  string
	Mandatory bool
	Immediate bool
}

var DefaultPublishOption = PublishOption{
	Exchange:  "",
	Mandatory: false,
	Immediate: false,
}
