package main

import (
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"demo/handler"
	"demo/lib/config"
	"demo/lib/httputil"
	"demo/setting"

	xormCore "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	initDependency()

	//in old go compiler, it is a must to enable multithread processing
	runtime.GOMAXPROCS(runtime.NumCPU())

	router := mux.NewRouter()
	uuidRegexp := `[[:alnum:]]{8}-[[:alnum:]]{4}-4[[:alnum:]]{3}-[89AaBb][[:alnum:]]{3}-[[:alnum:]]{12}`

	router.HandleFunc("/v1/cats/", handler.CatGetAll).Methods("GET")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", handler.CatGetOne).Methods("GET")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", handler.CatUpdate).Methods("PUT")
	router.HandleFunc("/v1/cats/{catId:"+uuidRegexp+"}", handler.CatDelete).Methods("DELETE")
	router.HandleFunc("/v1/cats/", handler.CatCreate).Methods("POST")

	http.Handle("/", router)
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}

// init the various object and inject the database object to the modules
func initDependency() {
	//the postgresql connection string
	connectStr := "host=" + config.GetStr(setting.DB_HOST) +
		" port=" + strconv.Itoa(config.GetInt(setting.DB_PORT)) +
		" dbname=" + config.GetStr(setting.DB_NAME) +
		" user=" + config.GetStr(setting.DB_USERNAME) +
		" password='" + config.GetStr(setting.DB_PASSWORD) + "'" +
		" sslmode=disable"

	db, err := xorm.NewEngine("postgres", connectStr)
	if err != nil {
		log.Panic("DB connection initialization failed", err)
	}

	db.SetMaxIdleConns(config.GetInt(setting.DB_MAX_IDLE_CONN))
	db.SetMaxOpenConns(config.GetInt(setting.DB_MAX_OPEN_CONN))
	db.SetColumnMapper(xormCore.SnakeMapper{})
	//uncomment it if you want to debug
	//db.ShowSQL = true
	//db.ShowErr = true

	httputil.Init(xormCore.SnakeMapper{})

	handler.Init(db)
}
