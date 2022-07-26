package pb

const (
	Success int32 = 200
	Error   int32 = 500

	NotFound         int32 = 404
	AccessDenied     int32 = 403
	WrongPassword    int32 = 5001
	InvalidOperation int32 = 5002
	DuplicateKey     int32 = 5003
)

const (
	MessageTransferTopic = "goim_message_transfer"
	MessageTransferGroup = "goim_group_transfer"
	MessagePersistGroup  = "goim_group_persist"

	MessagePushTopic = "goim_message_push"
)

const (
	UserServiceName    = "goim_rpc_service_user"
	FriendServiceName  = "goim_rpc_service_friend"
	GroupServiceName   = "goim_rpc_service_group"
	MessageServiceName = "goim_rpc_service_message"
	GatewayServiceName = "goim_rpc_service_gateway"
)
