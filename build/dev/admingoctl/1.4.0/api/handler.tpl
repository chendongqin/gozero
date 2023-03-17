package {{.PkgName}}

import (
    {{if .HasRequest}}"errors"{{end}}
	"net/http"
	{{if .HasRequest}}"reflect"{{end}}

    {{if .HasRequest}}"gozore-mall/common/errorx"
    "gozore-mall/common/utils"
    "strings"{{end}}
	"gozore-mall/common/response"
	{{.ImportPackages}}
    {{if .HasRequest}}
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    translations "github.com/go-playground/validator/v10/translations/zh"{{end}}
    "github.com/zeromicro/go-zero/rest/httpx"
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
            //兼容处理二次消费请求的body异常
            if strings.Index(err.Error(), "error: `EOF`") >= 0 {
                bodyErr := utils.GetReqJson(r.Context(), &req)
                if bodyErr == nil {
                    err = nil
                }
            }
            if err != nil {
                httpx.Error(w, errorx.NewHandlerError(errorx.ParamErrorCode, err.Error()))
                return
            }
        }

		validate := validator.New()
        validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
            name := fld.Tag.Get("label")
            return name
        })

        trans, _ := ut.New(zh.New()).GetTranslator("zh")
        validateErr := translations.RegisterDefaultTranslations(validate, trans)
        if validateErr = validate.StructCtx(r.Context(), req); validateErr != nil {
            for _, err := range validateErr.(validator.ValidationErrors) {
                httpx.Error(w, errorx.NewHandlerError(errorx.ParamErrorCode, errors.New(err.Translate(trans)).Error()))
                return
            }
        }

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			httpx.Error(w, err)
			return
		}

        {{if .HasResp}}response.Response(w, resp, err){{else}}response.Response(w, nil, err){{end}}
	}
}