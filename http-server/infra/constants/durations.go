package constants

import "time"

const REDIS_POOL_TIMEOUT time.Duration = 17 * time.Second
const REDIS_CONN_MAX_IDLE_TIME time.Duration = 5 * time.Minute
const REDIS_SESSION_TTL time.Duration = 48 * time.Hour
const REDIS_DESKTOP_AUTH_TTL time.Duration = 5 * time.Minute

const SERVER_READ_TIMEOUT time.Duration = 5 * time.Minute
const SERVER_WRITE_TIMEOUT time.Duration = 5 * time.Minute
const SERVER_IDLE_TIMEOUT time.Duration = 90 * time.Second

const AI_PROCCESS_TIMEOUT time.Duration = 61 * time.Second
const AI_WAIT_FREE time.Duration = 10 * time.Millisecond
const AI_WAIT_PREMIUM time.Duration = 2 * time.Millisecond

const EMAIL_LINK_TIMEOUT time.Duration = 24 * time.Hour
