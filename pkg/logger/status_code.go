package logger

type ResponseCode int

const (
	ResponseCodeInternalServer ResponseCode = 500
	ResponseCodeBadRequest     ResponseCode = 400
	ResponseCodeUnauthorized   ResponseCode = 401
	ResponseCodeForbidden      ResponseCode = 403
	ResponseCodeNotFound       ResponseCode = 404
	ResponseCodeNoContent      ResponseCode = 204
	ResponseCodeOK             ResponseCode = 200
	ResponseCodeCreated        ResponseCode = 201
	ResponseCodeAccepted       ResponseCode = 202
)

type ResponseMessage string

const (
	ResponseMessageInternalServer ResponseMessage = "internal server error"
	ResponseMessageBadRequest     ResponseMessage = "bad request"
	ResponseMessageUnauthorized   ResponseMessage = "unauthorized"
	ResponseMessageForbidden      ResponseMessage = "forbidden"
	ResponseMessageNotFound       ResponseMessage = "not found"
	ResponseMessageNoContent      ResponseMessage = "no content"
	ResponseMessageOK             ResponseMessage = "ok"
	ResponseMessageCreated        ResponseMessage = "item was created"
	ResponseMessageAccepted       ResponseMessage = "request accepted"
)

func ResponseCodeToMessage(code ResponseCode) ResponseMessage {
	switch code {
	case ResponseCodeInternalServer:
		return ResponseMessageInternalServer
	case ResponseCodeBadRequest:
		return ResponseMessageBadRequest
	case ResponseCodeUnauthorized:
		return ResponseMessageUnauthorized
	case ResponseCodeForbidden:
		return ResponseMessageForbidden
	case ResponseCodeNotFound:
		return ResponseMessageNotFound
	case ResponseCodeNoContent:
		return ResponseMessageNoContent
	case ResponseCodeOK:
		return ResponseMessageOK
	case ResponseCodeCreated:
		return ResponseMessageCreated
	case ResponseCodeAccepted:
		return ResponseMessageAccepted
	default:
		return ResponseMessageInternalServer
	}
}

func ResponseMessageToCode(message ResponseMessage) ResponseCode {
	switch message {
	case ResponseMessageInternalServer:
		return ResponseCodeInternalServer
	case ResponseMessageBadRequest:
		return ResponseCodeBadRequest
	case ResponseMessageUnauthorized:
		return ResponseCodeUnauthorized
	case ResponseMessageForbidden:
		return ResponseCodeForbidden
	case ResponseMessageNotFound:
		return ResponseCodeNotFound
	case ResponseMessageNoContent:
		return ResponseCodeNoContent
	case ResponseMessageOK:
		return ResponseCodeOK
	case ResponseMessageCreated:
		return ResponseCodeCreated
	case ResponseMessageAccepted:
		return ResponseCodeAccepted
	default:
		return ResponseCodeInternalServer
	}
}
