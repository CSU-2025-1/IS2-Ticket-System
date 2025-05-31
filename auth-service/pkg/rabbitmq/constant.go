package rabbitmq

const (
	ContentTypeProtobuf = "application/x-protobuf"
	ContentTypeJson     = "application/json"
)

type Mode int

const (
	ModeProto Mode = iota
	ModeJson
)

var (
	ModeToContentType = map[Mode]string{
		ModeProto: ContentTypeProtobuf,
		ModeJson:  ContentTypeJson,
	}
)
