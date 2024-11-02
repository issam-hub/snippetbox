package models

type TemplateData interface {
	IsAuth() bool
	GetCSRFToken() string
}

// CrossTemplates implements TemplateData
type CrossTemplates struct {
	IsAuthenticated bool
	CSRFToken       string
}

func (ct *CrossTemplates) IsAuth() bool {
	return ct.IsAuthenticated
}

func (ct *CrossTemplates) GetCSRFToken() string {
	return ct.CSRFToken
}

type contextKey string

var ContextClass = contextKey("templateData")

type isAuthKey string

var IsAuthenticatedContextKey = contextKey("isAuth")
