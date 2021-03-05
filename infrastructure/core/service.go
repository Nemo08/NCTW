package core

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

type Service struct {
	handlers []commandHandlerStruct
}

func (svc *Service) NewCommandHandler(commandName string, usecase func(sc ServiceContext) error, dataType interface{}) *commandHandlerStruct {
	var (
		chs commandHandlerStruct = commandHandlerStruct{}
	)

	chs.name = commandName
	chs.usecase = usecase
	chs.dataType = reflect.TypeOf(dataType)
	svc.handlers = append(svc.handlers, chs)
	return &chs
}

func (svc *Service) RunCommand(sc ServiceContext, name string) error {
	for _, v := range svc.handlers {
		if v.name == name {
			//валидация входящих данных
			err := svc.DataValidate(sc, v)
			if err != nil {
				return err
			}
			//выполняем юзкейс
			fmt.Println("get command", name)
			return v.usecase(sc)
		}
	}
	return errors.New("command not found")
}

func (svc *Service) getDataType(name string) interface{} {
	for _, v := range svc.handlers {
		if v.name == name {
			return v.dataType
		}
	}
	return nil
}

func (svc *Service) DataValidate(sc ServiceContext, comm commandHandlerStruct) error {
	return dataValidate(sc, comm)
}

func NewService() Service {
	return Service{}
}

func (svc *Service) EchoEndpoint(name string) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("X-Total-Count", "0")

		tp := svc.getDataType(name)
		spew.Dump("TP", tp)
		if err := ctx.Bind(&tp); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
		}

		l := newLogger()
		reqID := ctx.Response().Header().Get("X-Request-ID")
		if reqID != "" {
			l = *l.WithField("request_id", reqID)
		}

		customContext := ServiceContext{
			RequestData: tp,
			Log:         &l,
		}

		err := svc.RunCommand(customContext, name)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, customContext.ResponseData)
	}
}
