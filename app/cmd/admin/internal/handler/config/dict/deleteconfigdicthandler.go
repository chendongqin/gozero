package dict

import (
	"errors"
	"net/http"
	"reflect"

	"gozore-mall/app/cmd/admin/internal/logic/config/dict"
	"gozore-mall/app/cmd/admin/internal/svc"
	"gozore-mall/app/cmd/admin/internal/types"
	"gozore-mall/common/errorx"
	"gozore-mall/common/response"
	"gozore-mall/common/utils"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteConfigDictHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteConfigDictReq
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

		l := dict.NewDeleteConfigDictLogic(r.Context(), svcCtx)
		err := l.DeleteConfigDict(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		response.Response(w, nil, err)
	}
}
