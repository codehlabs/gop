# Generic Profile

Generic user/profile sign up and implemetation for User/Profile base applications. Takes care of generic user CRUD using up to date security standards.


## QUICK START
Quick usage example

```go
import (
	"net/http"
	"github.com/racg0092/gop"
)


func main() {

  gop.SetDriverConfig(DriverConfig{
    Conn: os.Getenv("db"),
    Database: "db",
    Collection: "users",
  )

  // using mongo db as the driver to save data in this case
  _, e := NewDriver(driver.MONGO, GetDriverConfig())
  if e != nil {
    panic(e)
  }

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


  // Saved data to the database
  if e := user.Save(GetDriver()); e != nil {
    // Handle error
  }

}
```
