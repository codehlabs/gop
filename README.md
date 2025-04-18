# User/Profile

Generic user sign up and implemetation for User/Profile bases applications. Takes care of generic user CRUD using up to date security standards.


## QUICK START
Sample of usage

```go
import (
	"net/http"
	"github.com/racg0092/gop"
	"github.com/racg0092/gop/driver"
)


func main() {
  mux := http.NewServerMux()
  //  Sign up or sign Path
  mux.HandleFunc("POST /auth", signup)
  server := http.Server{Addr: "127.0.0.1:8080", Handler: mux}
  if e := server.ListenAndServer(); e != nil {
    panic(e)
  }
}


func signup(w http.ResponseWriter, r *http.Request) {

  // Parsing the from data from your html form
  if e := r.ParseForm(); e != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // Returns an user struct
  user, e := gop.UserFromForm(r)
  if e != nil {
    // Handle error
  }

  // logic configuration. this is optional there are built in defaults
	gop.SetConfig(
    gop.Config{
      UniqueIDLength: 32,
      UseBuiltInSaveLogic: true
    }
  )

  // driver configuration required for saving to database
	config := driver.InitConfig{
    Conn: os.Getenv("db"),
    Database: "lobos",
    Collection: "users"
  }

  // using mongo db as the driver to save data in this case
  drv, e := driver.New(driver.MONGO, config)
  if e != nil {
    //Handle error
  }


  // Saved data to the database
  if e := user.Save(drv); e != nil {
    // Handle error
  }

}
```
