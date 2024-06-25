package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

type DB_CFG struct {
	Name     string
	Type     string
	Addr     string
	Port     string
	Login    string
	Password string
	DB       string
}
type Project struct {
	name       string
	group      string
	status     string
	script     string
	comments   string
	rate       string
	burst      string
	updatetime string
}
type PageData struct {
	Title   string
	Content string
}
type TO_SEND struct {
	Name string
	Type string
	Addr string
}
type PRE_SPDR_CFG struct {
	Name, Addr                  string
	TaskDB, ResultDB, ProjectDB TO_SEND
}
type SPDR_CFG_TO_SEND struct {
	Taskdb    string               `json:"taskdb"`
	Projectdb string               `json:"projectdb"`
	Resultdb  string               `json:"resultdb"`
	WebUI     PORT_SET_DATA        `json:"webui"`
	Fetcher   XMLRPC_PORT_SET_DATA `json:"fetcher"`
	Puppeteer PORT_SET_DATA        `json:"puppeteer"`
	Scheduler XMLRPC_PORT_SET_DATA `json:"scheduler"`
}
type SPDR_CFG struct {
	Addr   string
	Config SPDR_CFG_TO_SEND
}
type DEL_MSG struct {
	To_del string
}
type STATUS struct {
	Addr, Status string
	TSK          STATUS_DB
	RES          STATUS_DB
	PROJ         STATUS_DB
}
type STATUS_DB struct {
	Type, Addr, Status string
}
type RESALT_DATA struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}
type PORT_SET_DATA struct {
	Port int `json:"port"`
}
type XMLRPC_PORT_SET_DATA struct {
	Port int `json:"xmlrpc-port"`
}
type RESULT_FILTER struct {
	Source  string `json:"source_filter"`
	Key     string `json:"key_filter"`
	Content string `json:"content_filter"`
}

//	type RES_DB struct{
//		Type, Addr, Status string
//	}
//
//	type PROJ_DB struct{
//		Type, Addr, Status string
//	}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var sys_path = "."

// Обработчик статусов веб сокет
func wb_status_hadler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	go status_cheker(ws)
	for {
		if _, _, err := ws.NextReader(); err != nil {
			ws.Close()
			break
		}
	}
	log.Println("Client Connected")
	reader(ws)
}
func status_cheker(ws *websocket.Conn) {
	var cfg map[string]SPDR_CFG
	file, err := os.ReadFile(sys_path + "/config/cfg_Pyspyder.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}
	for {
		results := make(chan STATUS, len(cfg))
		for _, configItem := range cfg {
			go func(cfgItem SPDR_CFG) {
				results <- check_status(cfgItem)
			}(configItem)
		}
		for i := 0; i < len(cfg); i++ {
			req := <-results
			msg, err := json.Marshal(req)
			if err != nil {
				panic(err)
			}
			if req.Status == "offline" {
				if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("write:", err)
					return
				}
			}
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write:", err)
				return
			}
		}
		close(results)
		time.Sleep(5 * time.Second)
	}
}

//func result_cheker(ws *websocket.Conn) {
//	var cfg map[string]DB_CFG
//	file, err := os.ReadFile(sys_path + "/config/cfg_DB.json")
//	if err != nil {
//		panic(err)
//	}
//	err = json.Unmarshal(file, &cfg)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(len(cfg))
//	for {
//		results := make(chan []string, len(cfg))
//		var wg sync.WaitGroup
//
//		for _, configItem := range cfg {
//			if configItem.DB == "resultdb" {
//				wg.Add(1)
//				go func(cfgItem DB_CFG) {
//					defer wg.Done()
//					results <- get_results(cfgItem)
//				}(configItem)
//			}
//		}
//
//		go func() {
///			wg.Wait()
//		close(results)
//		}()
//
//		for req := range results {
//			msg, err := json.Marshal(req)
//			if err != nil {
//				log.Println("Ошибка при маршалинге:", err)
//				continue
//			}
//			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
//				log.Println("Ошибка записи в WebSocket:", err)
//			}
//		}
//		time.Sleep(5 * time.Second) // Задержка перед следующим циклом
//	}
//}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received: %sn", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// Маршрутизация главной страницы
func main_page(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("%v", r.URL.String())
	if r.Method != "GET" {
		fmt.Fprintln(w, "Тебе сюда нельзя!")
		return
	}
	if r.URL.String() == "/" {
		tmpl := template.Must(template.ParseFiles(sys_path+"/public/html/main_page.html", sys_path+"/public/templates/header.html"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, "")

	} else if r.URL.String() == "/?huy1=zalupa" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Jopa"))

	}
}

// Маршрутизация для конфигурации базы данных
func bd_settings(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v", r.URL.String())
	if r.URL.String() == "/bd_settings" {
		tmpl := template.Must(template.ParseFiles(sys_path+"/public/html/bd_settings.html", sys_path+"/public/templates/header.html"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, "")
	} else if r.URL.String() == "/bd_settings?bd=get_conf" {
		var m map[string]map[string]DB_CFG
		file, err := os.ReadFile(sys_path + "/config/cfg_DB.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &m)
		w.Header().Set("Content-Type", "application/json")
		w.Write(file)
	} else if r.URL.String() == "/bd_settings?bd=add_conf" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var req DB_CFG
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		fmt.Println(req)
		file, err := os.ReadFile(sys_path + "/config/cfg_DB.json")
		if err != nil {
			panic(err)
		}
		var cfg map[string]DB_CFG
		json.Unmarshal(file, &cfg)
		cfg[strconv.Itoa(len(cfg))] = req
		to_write, err := json.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		os.WriteFile(sys_path+"/config/cfg_DB.json", to_write, 0644)
		w.Write([]byte("succes"))

	} else if r.URL.String() == "/bd_settings?bd=del_conf" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var req DEL_MSG
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		var cfg_file map[string]DB_CFG
		file, err := os.ReadFile(sys_path + "/config/cfg_DB.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(file, &cfg_file)
		if err != nil {
			panic(err)
		}
		fmt.Println(cfg_file)
		delete(cfg_file, req.To_del)
		to_wrt, err := json.Marshal(cfg_file)
		if err != nil {
			panic(err)
		}
		os.WriteFile(sys_path+"/config/cfg_DB.json", to_wrt, 0644)
		w.Write([]byte("succes"))
	}
}

// Маршрутизация управления pyspyder
func pspdr_settings(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v", r.URL.String())
	if r.URL.String() == "/spyder_settings" {
		tmpl := template.Must(template.ParseFiles(sys_path+"/public/html/pspdr_settings.html", sys_path+"/public/templates/header.html"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, "")
		//Добавление конфигурации pyspyder
	} else if r.URL.String() == "/spyder_settings?spdr=add" {
		var req PRE_SPDR_CFG
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		var cfg map[string]SPDR_CFG
		file, err := os.ReadFile(sys_path + "/config/cfg_Pyspyder.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &cfg)
		addr_parts := strings.Split(req.Addr, ":")
		port_int, err := strconv.Atoi(addr_parts[1])
		if err != nil {
			panic(err)
		}
		cfg[req.Name] = SPDR_CFG{
			Addr: req.Addr,
			Config: SPDR_CFG_TO_SEND{
				Taskdb:    req.TaskDB.Type + "+taskdb://" + req.TaskDB.Addr + "/taskdb",
				Projectdb: req.ProjectDB.Type + "+projectdb://" + req.ProjectDB.Addr + "/projectdb",
				Resultdb:  req.ResultDB.Type + "+resultdb://" + req.ResultDB.Addr + "/resultdb",
				WebUI: PORT_SET_DATA{
					Port: port_int,
				},
				Fetcher: XMLRPC_PORT_SET_DATA{
					Port: port_int + 1,
				},
				Puppeteer: PORT_SET_DATA{
					Port: port_int + 2,
				},
				Scheduler: XMLRPC_PORT_SET_DATA{
					Port: port_int + 3,
				},
			},
		}
		to_wrt, err := json.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		os.WriteFile(sys_path+"/config/cfg_Pyspyder.json", to_wrt, 0644)
		w.Write([]byte("succes"))
		//Информация о доступных СУБД для конфигурирования pyspyder
	} else if r.URL.String() == "/spyder_settings?spdr=gt_opt" {
		cfg := read_dbs_cfg(sys_path + "/config/cfg_DB.json")
		cfg_list, err := json.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(cfg_list)
		// Вывод информации о pyspyder
	} else if r.URL.String() == "/spyder_settings?spdr=ref" {
		cfg_spdrs, err := os.ReadFile(sys_path + "/config/cfg_Pyspyder.json")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(cfg_spdrs)
		// Удаление конфигурации pyspyder
	} else if r.URL.String() == "/spyder_settings?spdr=del" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var req DEL_MSG
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		var cfg_file map[string]SPDR_CFG
		file, err := os.ReadFile(sys_path + "/config/cfg_Pyspyder.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(file, &cfg_file)
		if err != nil {
			panic(err)
		}
		delete(cfg_file, req.To_del)
		to_wrt, err := json.Marshal(cfg_file)
		if err != nil {
			panic(err)
		}
		os.WriteFile(sys_path+"/config/cfg_Pyspyder.json", to_wrt, 0644)
		w.Write([]byte("succes"))
	}
}

// Маршрутизация с результатами работы pyspyder
func result_page(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v", r.URL.String())
	if r.URL.String() == "/result" {
		tmpl := template.Must(template.ParseFiles(sys_path+"/public/html/result.html", sys_path+"/public/templates/header.html"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, "")
	} else if r.URL.String() == "/result?result=get" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var req RESULT_FILTER
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		fmt.Println(req)
		var cfg map[string]DB_CFG
		file, err := os.ReadFile(sys_path + "/config/cfg_DB.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(file, &cfg)
		if err != nil {
			panic(err)
		}
		results := make(chan []string, len(cfg))
		var wg sync.WaitGroup

		for _, configItem := range cfg {
			if configItem.DB == "resultdb" {
				wg.Add(1)
				go func(cfgItem DB_CFG) {
					defer wg.Done()
					results <- get_results(cfgItem, req)
				}(configItem)
			}
		}
		go func() {
			wg.Wait()
			close(results)
		}()
		var full_responce string
		for req := range results {
			full_responce += strings.Join(req, "")
			//fmt.Println(full_responce)
		}
		msg, err := json.Marshal(full_responce)
		if err != nil {
			log.Println("Ошибка при маршалинге:", err)
		}
		w.Write(msg)
	}

}

// Обмен данных с pyspyder-ами
func com_spyder(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v", r.URL.String())
	// Загрузка конфигурации в pyspyder
	if r.URL.String() == "/spyders?spdr=start" {
		ip, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Ошибка при получении IP адреса: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println(ip, port)
		var cfgs_spdr map[string]SPDR_CFG
		file, err := os.ReadFile(sys_path + "/config/cfg_Pyspyder.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(file, &cfgs_spdr)
		if err != nil {
			panic(err)
		}
		for t := range cfgs_spdr {
			parts := strings.Split(cfgs_spdr[t].Addr, ":")
			if ip == parts[0] {
				to_snd, err := json.Marshal(cfgs_spdr[t].Config)
				if err != nil {
					panic(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(to_snd)
				break
			}
		}
	}
}

// Список доступных СУБД для pyspyder
func read_dbs_cfg(path string) map[string]map[string]TO_SEND {
	var tmp map[string]DB_CFG
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &tmp)
	to_snd := make(map[string]map[string]TO_SEND)
	to_snd["resultdb"] = make(map[string]TO_SEND)
	to_snd["taskdb"] = make(map[string]TO_SEND)
	to_snd["projectdb"] = make(map[string]TO_SEND)
	cnt_0 := 0
	cnt_1 := 0
	cnt_2 := 0
	//Подготовка пакета с информацией о СУБД
	for i := range tmp {
		if tmp[i].DB == "resultdb" {
			if tmp[i].Login != "" && tmp[i].Password != "" {
				to_snd["resultdb"][strconv.Itoa(cnt_0)] = TO_SEND{
					Name: tmp[i].Name,
					Type: tmp[i].Type,
					Addr: tmp[i].Login + ":" + tmp[i].Password +
						"@" + tmp[i].Addr + ":" + tmp[i].Port,
				}
			} else {
				to_snd["resultdb"][strconv.Itoa(cnt_0)] = TO_SEND{
					Name: tmp[i].Name,
					Type: tmp[i].Type,
					Addr: tmp[i].Addr + ":" + tmp[i].Port,
				}
			}
			cnt_0++
		} else if tmp[i].DB == "taskdb" {
			if tmp[i].Login != "" && tmp[i].Password != "" {
				to_snd["taskdb"][strconv.Itoa(cnt_1)] = TO_SEND{
					Name: tmp[i].Name,
					Type: tmp[i].Type,
					Addr: tmp[i].Login + ":" + tmp[i].Password +
						"@" + tmp[i].Addr + ":" + tmp[i].Port,
				}
			} else {
				to_snd["taskdb"][strconv.Itoa(cnt_1)] = TO_SEND{
					Name: tmp[i].Name,
					Type: tmp[i].Type,
					Addr: tmp[i].Addr + ":" + tmp[i].Port,
				}
			}
			cnt_1++
		} else if tmp[i].DB == "projectdb" {
			if tmp[i].Login != "" && tmp[i].Password != "" {
				to_snd["projectdb"][strconv.Itoa(cnt_2)] = TO_SEND{
					Name: tmp[i].Name,
					Type: tmp[i].Type,
					Addr: tmp[i].Login + ":" + tmp[i].Password +
						"@" + tmp[i].Addr + ":" + tmp[i].Port,
				}
			} else {
				to_snd["projectdb"][strconv.Itoa(cnt_2)] = TO_SEND{
					Name: tmp[i].Name,
					Type: tmp[i].Type,
					Addr: tmp[i].Addr + ":" + tmp[i].Port,
				}
			}
			cnt_2++
		}
	}
	return to_snd
}

// Проверка статуса pyspyder online/offline + подключенные СУБД
func check_status(cfgObj SPDR_CFG) STATUS {
	task_db := make(chan STATUS_DB, 1)
	res_db := make(chan STATUS_DB, 1)
	proj_db := make(chan STATUS_DB, 1)
	url := "http://" + cfgObj.Addr
	go func() {
		task_db <- bd_status_check(cfgObj.Config.Taskdb, "+taskdb")
	}()
	go func() {
		res_db <- bd_status_check(cfgObj.Config.Resultdb, "+resultdb")
	}()
	go func() {
		proj_db <- bd_status_check(cfgObj.Config.Projectdb, "+projectdb")
	}()
	timeout := time.Duration(3 * time.Second) // Устанавливаем таймаут
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url) // Отправляем HTTP GET запрос
	if err != nil {
		// fmt.Printf("Ошибка при попытке 'ping': %sn", err)
		return STATUS{
			Addr:   cfgObj.Addr,
			Status: "offline 🔴",
			TSK:    <-task_db,
			RES:    <-res_db,
			PROJ:   <-proj_db,
		}
	}
	defer resp.Body.Close() // Не забываем закрыть тело ответа
	result := STATUS{
		Addr:   cfgObj.Addr,
		Status: "online 🟢",
		TSK:    <-task_db,
		RES:    <-res_db,
		PROJ:   <-proj_db,
	}
	close(task_db)
	close(res_db)
	close(proj_db)
	return result
}

// проверка статуса СУБД
func bd_status_check(db_url string, mask string) STATUS_DB {
	url := strings.Replace(db_url, mask, "", -1)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if strings.HasPrefix(url, "mongodb") {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				log.Fatal(err)
			}
		}()
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		url_parts := strings.Split(url, "/")
		if err = client.Ping(ctx, nil); err != nil {
			// fmt.Println("MongoDB недоступен:", err)
			url_parts := strings.Split(url, "/")
			return STATUS_DB{
				Type:   "mongodb",
				Addr:   url_parts[2] + "/" + url_parts[3],
				Status: "offline 🔴",
			}
		}
		return STATUS_DB{
			Type:   "mongodb",
			Addr:   url_parts[2] + "/" + url_parts[3],
			Status: "online 🟢",
		}
	}
	return STATUS_DB{}
}

// Загрузка результатов работы
func get_results(cfg DB_CFG, filter RESULT_FILTER) []string {
	// Устанавливаем контекст исполнения
	var divs []string
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var con_url string
	// Подключаемся к MongoDB
	if cfg.Login != "" && cfg.Password != "" {
		con_url = "mongodb://" + cfg.Login + ":" + cfg.Password + "@" + cfg.Addr + ":" + cfg.Port
	} else {
		con_url = "mongodb://" + cfg.Addr + ":" + cfg.Port
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(con_url))
	if err != nil {
		log.Println(err)
		to_send_err_str := "<div data-tmp=\"" + cfg.Name + "\"><h3>" + cfg.Addr +
			":" + cfg.Port + "  offline 🔴</h3></div>"
		divs = append(divs, to_send_err_str)
		return divs
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Ошибка при отключении от MongoDB:", err)
		}
	}()
	// Выбираем базу данных и коллекцию
	database := client.Database(cfg.DB)

	// Выполняем запрос Find для получения ключа result из всех документов
	collections, err := database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Println(err)
		to_send_err_str := "<div data-tmp=\"" + cfg.Name + "\"><h3>" + cfg.Addr + ":" +
			cfg.Port + "  offline 🔴</h3></div>"
		divs = append(divs, to_send_err_str)
		return divs
	}
	var srch_fltr bson.M
	if filter.Source == "" && filter.Key != "" {
		srch_fltr = bson.M{"result": bson.M{"$regex": "\"" + filter.Key + "\":", "$options": "i"}}
	} else if filter.Key == "" && filter.Source != "" {
		srch_fltr = bson.M{"url": bson.M{"$regex": filter.Source, "$options": "i"}}
	} else if filter.Source != "" && filter.Key != "" {
		srch_fltr = bson.M{
			"$and": []bson.M{
				{"url": bson.M{"$regex": filter.Source, "$options": "i"}},
				{"result": bson.M{"$regex": "\"" + filter.Key + "\":", "$options": "i"}},
			},
		}
	} else {
		srch_fltr = bson.M{}
	}
	for _, collectionName := range collections {
		collection := database.Collection(collectionName)
		cur, err := collection.Find(ctx, srch_fltr, options.Find().SetProjection(bson.D{{"result", 1}, {"url", 1}}))
		if err != nil {
			log.Println(err)
			to_send_err_str := "<div data-tmp=\"" + cfg.Name + "\"><h3>" + cfg.Addr + ":" +
				cfg.Port + "  offline 🔴</h3></div>"
			divs = append(divs, to_send_err_str)
			return divs
		}
		defer cur.Close(ctx)
		// Итерируемся по курсору и читаем документы
		for cur.Next(ctx) {
			var result map[string]interface{}
			err := cur.Decode(&result)
			if err != nil {
				log.Println(err)
				continue
			}
			resultRaw, ok := result["result"].(string)
			if !ok {
				log.Println("Не удалось получить result или result не строка")
				continue
			}
			link_row, ok := result["url"].(string)
			if !ok {
				log.Println("Не удалось получить result или result не строка")
				continue
			}
			var resultData map[string]interface{}
			if err := json.Unmarshal([]byte(resultRaw), &resultData); err != nil {
				log.Println("Ошибка при десериализации result:", err)
				continue
			}
			source := strings.Split(link_row, "/")
			out_chek := false
			var to_send_str string
			//to_send_str := "<div data-tmp=\"" + cfg.Name + "\" class=\"" + source[2] + "\"><a href=\"" + link_row +
			//	"\">" + source[2] + "</a>"
			for i := range resultData {
				tmp, ok := resultData[i].(string)
				if !ok {
					log.Println("Не удалось получить result или result не строка")
					continue
				}
				if filter.Content != "" {
					if strings.Contains(strings.ToLower(tmp), strings.ToLower(filter.Content)) {
						if !out_chek {
							to_send_str += "<div data-tmp=\"" + cfg.Name + "\" class=\"" + source[2] +
								"\"><a href=\"" + link_row + "\">" + source[2] + "</a>"
							out_chek = true
						}
						to_send_str += "<div><p class=\"" + i + "\">" + i + " : " + tmp + "</p></div>"
					}
				} else {
					if !out_chek {
						to_send_str += "<div data-tmp=\"" + cfg.Name + "\" class=\"" + source[2] +
							"\"><a href=\"" + link_row + "\">" + source[2] + "</a>"
						out_chek = true
					}
					to_send_str += "<div><p class=\"" + i + "\">" + i + " : " + tmp + "</p></div>"
				}
			}
			if out_chek {
				to_send_str += "</div>"
				divs = append(divs, to_send_str)
			}
		}
		if err := cur.Err(); err != nil {
			log.Println(err)
			to_send_err_str := "<div data-tmp=\"" + cfg.Name + "\"><h3>" + cfg.Addr +
				":" + cfg.Port + "  offline 🔴</h3></div>"
			divs = append(divs, to_send_err_str)
			return divs
		}
	}
	return divs
}
func main() {
	pnt_1 := time.Now()
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", main_page)
	http.HandleFunc("/bd_settings", bd_settings)
	http.HandleFunc("/spyder_settings", pspdr_settings)
	http.HandleFunc("/spyders", com_spyder)
	http.HandleFunc("/result", result_page)
	http.HandleFunc("/ws", wb_status_hadler)
	// http.HandleFunc("/ws_res", wb_result_hadler)
	pnt_2 := time.Now()
	fmt.Println(pnt_2.Sub(pnt_1))
	http.ListenAndServe(":8000", nil)

}
