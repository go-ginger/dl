package query

import gm "github.com/go-ginger/models"

type IParseResult interface {
	GetQuery() (query interface{})
	GetSort() (sort interface{})
	GetParams() (params []interface{})
}

type IParser interface {
	Parse(request gm.IRequest) (result IParseResult)
}
