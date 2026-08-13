package httperr

import "errors"

var (
	Err_DuplicateEmail          = errors.New("error.duplicate-email")
	Err_UserNotFound            = errors.New("error.user-not-found")
	Err_ContextDeadlineExceeded = errors.New("error.context-deadline-exceeded")
	Err_ContextCanceled         = errors.New("error.context-canceled")
	Err_UniqueViolation         = errors.New("error.unique-violation")
	Err_DbTimeout               = errors.New("error.db-timeout")
	Err_DbNetworkTemporary      = errors.New("error.db-network-temporary")
	Err_DbNetwork               = errors.New("error.db-network")
	Err_NotDeleted              = errors.New("error.not-deleted")
	Err_NotUpdated              = errors.New("error.not-updated")
)

var (
	Err_RedisKeyNotFound      = errors.New("error.redis-key-not-found")
	Err_RedisTimeout          = errors.New("error.redis-timeout")
	Err_RedisDeadlineExceeded = errors.New("error.redis-deadline-exceeded")
	Err_RedisCanceled         = errors.New("error.redis-canceled")
	Err_RedisNetwork          = errors.New("error.redis-network")
	Err_RedisOperationFailed  = errors.New("error.redis-operation-failed")
)
