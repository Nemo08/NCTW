package core

type CommandOption func(*commandHandlerStruct)

//Name добавляет имя команды
func Name(name string) CommandOption {
	return func(chs *commandHandlerStruct) {
		chs.name = name
	}
}

//Description добавляет описание команды
func Description(descr string) CommandOption {
	return func(chs *commandHandlerStruct) {
		chs.description = descr
	}
}

//DataType принимает тип данных, с которым будет работать команда
//Также в типе должны присутствовать аннотации валидатора
func DataType(dataStruct interface{}) CommandOption {
	return func(chs *commandHandlerStruct) {
		chs.dataType = dataStruct
	}
}
