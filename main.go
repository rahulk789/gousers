package main
import (
    "log"
    "fmt"
    "github.com/gorilla/mux"
    "net/http"
    "strings"
    "github.com/go-sql-driver/mysql"
    "os"
    "database/sql" 
)
type USERS struct {
    name string
    pass string
}
var db *sql.DB
func main() {
    // Capture connection properties.
    cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "users",
        AllowNativePasswords : true,
    }
    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")
    UserS, err := connectTheUser("new user")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("users found: %v\n", UserS)
    //server handling
    r :=mux.NewRouter()
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))
    r.HandleFunc("/",HomeHandler)
    r.HandleFunc("/website",WebsiteHandler)
    log.Fatal(http.ListenAndServe(":8080",r))
}
//DB functions
func connectTheUser(x string) ([]USERS, error) {
    var Users []USERS
    rows, err := db.Query("SELECT * FROM USERS")
    if err != nil {
        return nil, fmt.Errorf("connectTheUser %q: %v", x, err)
   }
    defer rows.Close()
    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var usr USERS
        if err := rows.Scan(&usr.name, &usr.pass,); err != nil {
            return nil, fmt.Errorf("connectTheUser %q: %v", x, err)
        }
        Users = append(Users, usr)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("connectTheUser %q: %v", x, err)
    }
    return Users, nil
}

func addUser(usr USERS) ( error) {
    _,err :=  db.Exec("INSERT INTO USERS (name, pass) VALUES (?, ?)", usr.name, usr.pass)
    if err != nil {
        return fmt.Errorf("addUSERS: %v", err)
    }
    return nil
}
//server handling functions
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./static/index.html")
}
func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
            r.ParseForm()
            stringByte1 := strings.Join(r.Form["username"]," ")
            stringByte2 := strings.Join(r.Form["password"]," ")
            fmt.Fprintf(w,"\nusername: ")
            fmt.Fprintf(w,stringByte1)
            fmt.Fprintf(w,"\npassword: ")
            fmt.Fprintf(w,stringByte2)
            addUser(USERS{
                name: strings.Join(r.Form["username"]," "),
                pass: strings.Join(r.Form["password"]," ")})
         }
