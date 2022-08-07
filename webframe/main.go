package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
	tmpl   *template.Template
}

type speedResult struct {
	ID            int64  `gorm:"column:id"`
	TimeStamp     string `gorm:"column:TimeStamp"`
	DownloadSpeed uint64 `gorm:"column:DownloadSpeed"`
	UploadSpeed   int64  `gorm:"column:UploadSpeed"`
	Latency       int64  `gorm:"column:Latency"`
	PublicIp      string `gorm:"column:PublicIp"`
	ISP           string `gorm:"column:ISP"`
	Peers         string `gorm:"column:Peers"`
}

type tabler interface {
	TableName() string
}

func (speedResult) TableName() string {
	return "speedtest"
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) Initialize(user, password, dbname, host string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbname)
	fmt.Println(connectionString)
	var err error
	a.DB, err = gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}

}

func (a *App) Run(addr string) {

	log.Fatal(http.ListenAndServe("localhost:8010", a.Router))
}

func getResults(db *gorm.DB) []speedResult {
	var speedtest []speedResult
	if err := db.Find(&speedtest).Error; err != nil {
		log.Fatal(err)
	}
	log.Printf("%d rows found.", len(speedtest))
	return speedtest
}

func (a *App) getResults(w http.ResponseWriter, r *http.Request) {
	results := getResults(a.DB)
	respondWithJSON(w, http.StatusOK, results)
}

func (a *App) indexFile(w http.ResponseWriter, r *http.Request) {
	results := getResults(a.DB)
	fmt.Println(results[0])
	fmt.Printf("t1 : %T\n", results[0])
	a.tmpl.ExecuteTemplate(w, "index.html", results)
}

//func (a *App) InitializeRoutes() {

//}

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
		os.Getenv("MYSQL_HOST"),
	)

	a.tmpl = template.Must(template.ParseGlob("templates/*.html"))
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", a.indexFile).Methods("GET")
	a.Router.HandleFunc("/results", a.getResults).Methods("GET")
	//a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))

	fs := http.FileServer(http.Dir("./static/"))
	a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", a.Router)
	a.Run("8010")

}
