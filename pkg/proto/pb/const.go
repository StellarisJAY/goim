package pb

const (
	Success int32 = 200
	Error   int32 = 500

	NotFound      int32 = 404
	AccessDenied  int32 = 403
	WrongPassword int32 = 5001
)

const (
	MessageTransferTopic = "goim_message_transfer"
	MessageTransferGroup = "goim_group_transfer"
	MessagePersistGroup  = "goim_group_persist"
)
