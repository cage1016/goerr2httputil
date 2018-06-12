package custom_error1

type ErrorCode uint

// 400
const (
	BroadcastErrorQueryBroadcastRepoFindBadRequest ErrorCode = iota + 4000500
	BroadcastErrorCreateBroadcastRepoCreateBadRequest
	BroadcastErrorDeleteBroadcastRepoDeleteBadRequest
)

// 500
const (
	BroadcastErrorQueryBroadcastRepoFindServerInternalError ErrorCode = iota + 5000500
	BroadcastErrorCreateBroadcastRepoCreateServerInternalError
	BroadcastErrorDeleteBroadcastRepoDeleteServerInternalError
)

