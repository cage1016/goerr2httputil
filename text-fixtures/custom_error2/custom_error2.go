package custom_error2

type ErrorCode uint

// 400
const (
	BroadcastLangErrorQueryBroadcastLangRepoFindBadRequest ErrorCode = iota + 4001200
)

// 500
const (
	BroadcastLangErrorQueryBroadcastLangRepoFindServerInternalError ErrorCode = iota + 5001200 // custom error message
)
