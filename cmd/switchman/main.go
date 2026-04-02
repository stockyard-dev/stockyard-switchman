package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-switchman/internal/server";"github.com/stockyard-dev/stockyard-switchman/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9970"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./switchman-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("switchman: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Switchman — blue-green deploy manager\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("switchman: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
