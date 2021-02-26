package core

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type service struct {
	handlers []commandHandlerStruct
}

func (svc *service) NewCommandHandler(commandName string, usecase func(sc ServiceContext) error, dataType interface{}) error {
	var (
		chs commandHandlerStruct = commandHandlerStruct{}
	)

	chs.name = commandName
	chs.usecase = usecase
	svc.handlers[len(svc.handlers)] = chs
	return nil
}

func (svc *service) RunCommand(sc ServiceContext, name string) error {
	for _, v := range svc.handlers {
		if v.name == name {
			//валидация входящих данных
			err := svc.DataValidate(sc, v)
			if err != nil {
				return err
			}
			//выполняем юзкейс
			return v.usecase(sc)
		}
	}
	return errors.New("command not found")
}

func (svc *service) getDataType(name string) interface{} {
	for _, v := range svc.handlers {
		if v.name == name {
			return v.dataType
		}
	}
	return nil
}

func (svc *service) EchoEndpoint(ctx echo.Context, command string) (err error) {
	ctx.Response().Header().Set("X-Total-Count", "0")

	tp := svc.getDataType(command)
	if err = ctx.Bind(tp); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	l := newLogger()
	reqID := ctx.Response().Header().Get("X-Request-ID")
	if reqID != "" {
		l = *l.WithField("request_id", reqID)
	}

	customContext := ServiceContext{
		requestData: tp,
		Log:         &l,
	}

	err = svc.RunCommand(customContext, command)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, customContext.responseData)
}

func (svc *service) DataValidate(sc ServiceContext, comm commandHandlerStruct) error {
	return dataValidate(sc, comm)
}

func NewService() service {
	return service{}
}
