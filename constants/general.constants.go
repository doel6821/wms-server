package constants

// SERVICE_CODE ..
const (
	SERVICE_CODE = "01"
	SERVICE_NAME = "wms-server"
)

const (
	REF_FILE = iota
	DURATION
	ROW_AFFECTED
	QUERY
)

const (
	LOG_OTHER = "OTHER_LOGS"
	LOG_QUERY = "QUERY_LOGS"
	START     = "Start"
)

const (
	PENDING       = "pending"
	ON_PROCESS    = "on process"
	ON_PACKING    = "on packing"
	ON_ALLOCATION = "on allocation"
	COMPLETE      = "complete"
)
