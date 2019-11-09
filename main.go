package main

import (
	"fmt"
	"strings"
)

type Location struct {
	Transitions []string
	Status      string
	Items       []string
}

var LocationMap = map[string]*Location{
	"коридор": {
		[]string{"кухня", "комната", "улица"},
		"ничего интересного",
		[]string{},
	},

	"комната": {
		[]string{"коридор"},
		"ты в своей комнате",
		[]string{"ключи", "конспекты", "рюкзак"},
	},

	"кухня": {
		[]string{"коридор"},
		"кухня, ничего интересного",
		[]string{"чай"},
	},

	"улица": {
		[]string{"домой"},
		"на улице весна",
		[]string{},
	},
}

// Состояние комнаты
type RoomStat struct {
	NameOfRoom string
	Items      []string
}

//Добавляем предметы в комнату
func (r *RoomStat) defolt(st string) {
	for k, val := range LocationMap {
		if k == st {
			r.Items = val.Items
		}
	}
}

//Удаляем предмет
func (r *RoomStat) takeItem(item string) bool {
	var idx int = -1
	for i, val := range r.Items {
		if item == val {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}
	copy(r.Items[idx:], r.Items[idx+1:]) //	Удаление предмета без угрозы
	r.Items[len(r.Items)-1] = ""         //	утечки памяти при инициализации
	r.Items = r.Items[:len(r.Items)-1]   //
	return true
}

//Возврщает состояние комнаты в виде строки
func (r RoomStat) getRoomStatus() string {
	var st2 string = ""
	for _, val := range r.Items {
		switch val {
		case "ключи":
			st2 += "на столе: " + val
		case "конспекты":
			if CheckKey {
				st2 += "на столе: " + val
			} else {
				st2 += ", " + val
			}
		case "рюкзак":
			st2 += ", на стуле: " + val
		case "чай":
			st2 += "на столе: " + val
		default:
			st2 = "пустая комната"
		}
	}
	return st2
}

//Инициализация
func (r *RoomStat) Init(m MyStat) {
	r.NameOfRoom = m.MyLoc
	r.defolt(r.NameOfRoom)
}

//Состояние данного игрока
type MyStat struct {
	MyLoc string   //Локация игрока
	Items []string //Предметы (подобранные)
}

//Инициализация
func (m *MyStat) Init() {
	for k := range LocationMap {
		m.MyLoc = k
		R.Init(*m)
	}
	m.MyLoc = "кухня"
	m.Items = []string{}
}

//Переход между комнатами
func (m *MyStat) changeMyLoc(r *RoomStat, toLoc string) bool {
	for _, val := range goTo(m.MyLoc) {
		if val == toLoc {
			m.MyLoc = toLoc
			r.Init(*m)
			return true
		}
	}
	return false
}

//Возвращает список комнат, в которые можно попасть из данной комнаты
func goTo(myLoc string) []string {
	for k, val := range LocationMap {
		if myLoc == k {
			return val.Transitions
		}
	}
	return []string{}
}

//Возврщает состояние игрока в виде строки
func (m *MyStat) getMyStatus() string {
	for k, val := range LocationMap {
		if m.MyLoc == k {
			return val.Status
		}
	}
	return ""
}

//Добавляем предмет к игроку
func (m *MyStat) take(r *RoomStat, item string) bool {
	if r.takeItem(item) {
		m.Items = append(m.Items, item)
		return true
	}
	return false
}

//Переход между комнатами (учитываем глобальное состояние комнаты R)
func (m *MyStat) Go(to string) string {
	//R.Init(*m) //-
	if m.changeMyLoc(&R, to) {
		for k, val := range LocationMap {
			if k == to {
				if to == "улица" {
					if CheckDoor {
						return val.Status + m.GoToStat()
					}
					m.MyLoc = "коридор"
					R.NameOfRoom = "коридор"
					return "дверь закрыта"

				}
				return val.Status + m.GoToStat()
			}
		}
	}
	return "нет пути в " + to
}

//Подбор предмета (учитываем глобальное состояние комнаты R)
func (m *MyStat) TakeStat(item string) string {
	//var st string = ""
	var idx = -1
	for i, val := range R.Items {
		if val == item {
			idx = i
		}
	}
	if R.takeItem(item) {
		if item == "рюкзак" {
			m.Items = append(m.Items, item)
			CheckBag = true
			return "вы надели: " + item
		} else {
			if CheckBag {
				m.Items = append(m.Items, item)
				if item == "ключи" {
					CheckKey = true
				}
				return "предмет добавлен в инвентарь: " + item
			}
			R.Items = append(R.Items, "")
			copy(R.Items[idx+1:], R.Items[idx:])
			R.Items[idx] = item
			return "некуда класть"
		}
	} else {
		return "нет такого"
	}
}

//Возвращает список комнат, в которые можно попасть из данной комнаты - R
func (m *MyStat) GoToStat() string {
	var st1 string = ". можно пройти -"
	for _, val := range goTo(m.MyLoc) {
		st1 = st1 + " " + val + ","
	}
	return st1
}

//Осмотр комнаты
func (m *MyStat) Look() string {
	if R.getRoomStatus() == "" {
		return "пустая комната" + m.GoToStat()
	}
	if R.NameOfRoom == "кухня" {
		if CheckBag {
			return "ты находишься на кухне, " + R.getRoomStatus() + ", надо идти в универ" + m.GoToStat()
		}
		return "ты находишься на кухне, " + R.getRoomStatus() + ", надо собрать рюкзак и идти в универ" + m.GoToStat()
	}
	return R.getRoomStatus() + m.GoToStat()
}

var Furn = map[string]string{
	"ключи":     "дверь",
	"конспекты": "шкаф",
}

//Применение предмета
func (m *MyStat) Apply(it string, furn string) string {
	var idx = -1
	//var st = ""
	for i, val := range m.Items {
		if it == val {
			idx = i
			break
		}
	}
	if idx == -1 {
		return "нет предмета в инвентаре - " + it
	}
	if Furn[it] != furn {
		return "не к чему применить"
	}
	if furn == "дверь" {
		CheckDoor = true
		return "дверь открыта"
	}
	return "конспекты в шкафу"
}

var R RoomStat
var CheckBag bool
var CheckDoor bool
var CheckKey bool
var CheckVisitRoom bool
var CheckVisitKitchen bool
var m MyStat

func main() {
	initGame()
	//1 - test
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("надеть рюкзак"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("взять конспекты"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("идти улица"))
	//2-test
	fmt.Println("---------------------------")
	initGame()
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("завтракать"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("надеть рюкзак"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("взять телефон"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять конспекты"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти кухня"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти улица"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("применить телефон шкаф"))
	fmt.Println(handleCommand("применить ключи шкаф"))
	fmt.Println(handleCommand("идти улица"))
}

func initGame() {
	/*
		эта функция инициализирует игровой мир - все команты
		если что-то было - оно корректно перезатирается
	*/
	//R.InitAllRooms()

	CheckBag = false
	CheckDoor = false
	CheckKey = false
	CheckVisitRoom = false
	CheckVisitKitchen = false
	m.Init()
	R.Init(m)
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/
	var sl = strings.Split(command, " ")
	switch sl[0] {
	case "идти":
		var answer = m.Go(sl[1])
		if len(answer) > 0 && answer[len(answer)-1] == ',' {
			answer = answer[:len(answer)-1]
		}
		if sl[1] == "комната" && !CheckVisitRoom && answer != "неизвестная команда" {
			R.Items = []string{"ключи", "конспекты", "рюкзак"}
			CheckVisitRoom = true
		} else if sl[1] == "кухня" && !CheckVisitKitchen && answer != "неизвестная команда" {
			R.Items = []string{"чай"}
			CheckVisitKitchen = true
		}
		return answer
	case "взять":
		var answer = m.TakeStat(sl[1])
		if len(answer) > 0 && answer[len(answer)-1] == ',' {
			answer = answer[:len(answer)-1]
		}
		return answer
	case "надеть":
		var answer = m.TakeStat(sl[1])
		if len(answer) > 0 && answer[len(answer)-1] == ',' {
			answer = answer[:len(answer)-1]
		}
		return answer
	case "применить":
		var answer = m.Apply(sl[1], sl[2])
		if len(answer) > 0 && answer[len(answer)-1] == ',' {
			answer = answer[:len(answer)-1]
		}
		return answer
	case "осмотреться":
		var answer = m.Look()
		if len(answer) > 0 && answer[len(answer)-1] == ',' {
			answer = answer[:len(answer)-1]
		}
		return answer
	default:
		return "неизвестная команда"
	}
}
