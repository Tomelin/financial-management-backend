package entity

import (
	"runtime"
	"strings"
)

const (
	ErrAlreadyExists = "already exists"
	ErrNotFound      = "not found"
	ErrRequired      = "is required"
	ErrCannotEmpty   = "cannot be empty"
)

type QueryClause string

const (
	QueryClauseAnd QueryClause = "and"
	QueryClauseOr  QueryClause = "or"
)

type QueryDBClause struct {
	Clause  QueryClause `json:"clause"`
	Queries []QueryDB   `json:"queries"`
}

type QueryDB struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Condition string `json:"condition"`
}

type QueryFirebase string

const (
	QueryFirebaseEqual              QueryFirebase = "=="
	QueryFirebaseNotEqual           QueryFirebase = "!="
	QueryFirebaseLessThan           QueryFirebase = "<"
	QueryFirebaseLessThanOrEqual    QueryFirebase = "<="
	QueryFirebaseGreaterThan        QueryFirebase = ">"
	QueryFirebaseGreaterThanOrEqual QueryFirebase = ">="
	QueryFirebaseArrayContains      QueryFirebase = "array-contains"
	QueryFirebaseArrayContainsAny   QueryFirebase = "array-contains-any"
	QueryFirebaseIn                 QueryFirebase = "in"
	QueryFirebaseNotIn              QueryFirebase = "not-in"
)

func QueryFirebaseString(q string) string {

	switch q {
	case "!=":
		return "!="
	case "<":
		return "<"
	case "<=":
		return "<="
	case ">":
		return ">"
	case ">=":
		return ">="
	case "array-contains":
		return "array-contains"
	case "array-contains-any":
		return "array-contains-any"
	case "in":
		return "in"
	case "not-in":
		return "not-in"
	default:
		return "=="
	}
}

type ModuleError struct {
	Module string       `json:"module"`
	Method string       `json:"method"`
	Layer  string       `json:"layer"`
	Code   ResponseCode `json:"code"`
	Err    string       `json:"message"`
}

type ApplicationLayer string

const (
	ApplicationLayerRepository ApplicationLayer = "Data Access Layer"
	ApplicationLayerService    ApplicationLayer = "Business Logic"
	ApplicationLayerHandler    ApplicationLayer = "Controllers"
	ApplicationLayerMiddleware ApplicationLayer = "Middleware"
	ApplicationLayerEntity     ApplicationLayer = "Model"
)

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

// Error
// Retorna um erro formatado com as informações do módulo, método, camada e código do erro
//
// Handlers (ou Controllers): Ponto de entrada da aplicação. Recebem requisições, chamam a camada de serviço e retornam respostas.
// Middleware: Interceptam requisições antes de chegarem aos handlers, executando lógica transversal (autenticação, autorização, logging, etc.).
// Service (ou Business Logic): Contém a lógica de negócio da aplicação. Orquestra as operações e interage com a camada de repositório.
// Repository (ou Data Access Layer): Responsável pela comunicação com o banco de dados ou outras fontes de dados. Abstrai a persistência de dados.
// Entity (ou Model): Representa os dados da aplicação (estruturas de dados, classes, etc.).
//
// Exemplo de uso:
// return nil, entity.Error(err.Error(), "wallet", "Create", "repository")
func Error(err, module, method string, layer ApplicationLayer, code ResponseCode) *ModuleError {

	m := ModuleError{
		Module: module,
		Method: method,
		Layer:  string(ApplicationLayer(layer)),
		Code:   code,
		Err:    err,
	}

	return &m
}

type ModuleLayer string

const (
	ModuleLayerUser          ModuleLayer = "user"
	ModuleLayerTenant        ModuleLayer = "tenant"
	ModuleLayerPlan          ModuleLayer = "plan"
	ModuleLayerAuthorization ModuleLayer = "authorization"
	ModuleLayerMiddleware    ModuleLayer = "middleware"
	ModuleLayerCategory      ModuleLayer = "transaction_category"
	ModuleLayerWallet        ModuleLayer = "wallet"
	ModuleLayerTransaction   ModuleLayer = "transaction"
)

func ErrorDebug(err string, code ResponseCode) *ModuleError {

	var funcParts []string = []string{"entity"}
	mErr := ModuleError{
		Code: code,
		Err:  err,
	}

	pc, file, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		funcParts = strings.Split(fn.Name(), "/")
		funcPart := funcParts[len(funcParts)-1]
		funcParts = strings.Split(funcPart, ".")
		method := funcParts[len(funcParts)-1]
		mErr.Method = method
	}

	fileParts := strings.Split(file, "/")
	filePart := fileParts[len(fileParts)-1]
	filePart = strings.Split(filePart, ".")[0]
	mErr.GetModule(&filePart)
	mErr.GetLayer(funcParts[0])

	return &mErr
}

func (m *ModuleError) GetModule(module *string) {
	var x string
	switch strings.ToLower(*module) {
	case "tenant":
		m.Module = string(ModuleLayerTenant)
	case "user":
		m.Module = string(ModuleLayerUser)
	case "transactioncategory":
		x = string(ModuleLayerCategory)
	case "plan":
		m.Module = string(ModuleLayerPlan)
	case "authorization":
		m.Module = string(ModuleLayerAuthorization)
	case "middleware":
		m.Module = string(ModuleLayerMiddleware)
	case "wallet":
		m.Module = string(ModuleLayerWallet)
	case "transaction":
		m.Module = string(ModuleLayerTransaction)
	default:
		m.Module = "unknown"
	}

	m.Module = x

}

func (m *ModuleError) GetLayer(layer string) {
	switch strings.ToLower(layer) {
	case "web":
		m.Layer = string(ApplicationLayerHandler)
	case "service":
		m.Layer = string(ApplicationLayerService)
	case "repository":
		m.Layer = string(ApplicationLayerRepository)
	case "middleware":
		m.Layer = string(ApplicationLayerMiddleware)
	case "entity":
		m.Layer = string(ApplicationLayerEntity)
	default:
		m.Layer = string(ApplicationLayerRepository)
	}
}
