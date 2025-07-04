package constants

// meta
const (
	RC_SUCCESS                   = "success"
	RC_BADREQUEST                = "badRequest"
	RC_INVALID_FORMAT            = "invalidFormat"
	RC_MANDATORY                 = "mandatory"
	RC_UNAUTHORIZED              = "unauthorized"
	RC_INVALID_TOKEN             = "invalidToken"
	RC_EXPIRED_TOKEN             = "TokenExpired"
	RC_INVALID_MERCHANT          = "invalidMerchant"
	RC_GENERAL_ERROR             = "general"
	RC_INVALID_EMAIL_OR_PASSWORD = "invalidEmailOrPassword"
	RC_EMAIL_ALREADY_USED        = "emailAlreadyUsed"
	RC_STORE_NAME_ALREADY_USED   = "storeNameAlreadyUsed"
	RC_PHONE_NUMBER_ALREADY_USED = "phoneNumberAlreadyUsed"
	RC_LOCATION_EXIST            = "locationExist"
	RC_TIMEOUT                   = "timeoutRequest"
)
const (
	HEAD_CONTENT_TYPE  = "Content-Type"
	HEAD_TIMESTAMP     = "X-TIMESTAMP"
	HEAD_CLIENTKEY     = "X-CLIENT-KEY"
	HEAD_SIGNATURE     = "X-SIGNATURE"
	HEAD_AUTHORIZATION = "Authorization"
	HEAD_PARTNERID     = "X-PARTNER-ID"
	HEAD_EXTERNALID    = "X-EXTERNAL-ID"
	HEAD_CHANNELID     = "CHANNEL-ID"
)
