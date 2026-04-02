package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Deployment struct {
	ID string `json:"id"`
	Service string `json:"service"`
	Version string `json:"version"`
	Environment string `json:"environment"`
	Strategy string `json:"strategy"`
	ActiveSlot string `json:"active_slot"`
	Status string `json:"status"`
	DeployedBy string `json:"deployed_by"`
	DeployedAt string `json:"deployed_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"switchman.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS deployments(id TEXT PRIMARY KEY,service TEXT NOT NULL,version TEXT DEFAULT '',environment TEXT DEFAULT 'production',strategy TEXT DEFAULT 'blue_green',active_slot TEXT DEFAULT 'blue',status TEXT DEFAULT 'pending',deployed_by TEXT DEFAULT '',deployed_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Deployment)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO deployments(id,service,version,environment,strategy,active_slot,status,deployed_by,deployed_at,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Service,e.Version,e.Environment,e.Strategy,e.ActiveSlot,e.Status,e.DeployedBy,e.DeployedAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Deployment{var e Deployment;if d.db.QueryRow(`SELECT id,service,version,environment,strategy,active_slot,status,deployed_by,deployed_at,created_at FROM deployments WHERE id=?`,id).Scan(&e.ID,&e.Service,&e.Version,&e.Environment,&e.Strategy,&e.ActiveSlot,&e.Status,&e.DeployedBy,&e.DeployedAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Deployment{rows,_:=d.db.Query(`SELECT id,service,version,environment,strategy,active_slot,status,deployed_by,deployed_at,created_at FROM deployments ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Deployment;for rows.Next(){var e Deployment;rows.Scan(&e.ID,&e.Service,&e.Version,&e.Environment,&e.Strategy,&e.ActiveSlot,&e.Status,&e.DeployedBy,&e.DeployedAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Deployment)error{_,err:=d.db.Exec(`UPDATE deployments SET service=?,version=?,environment=?,strategy=?,active_slot=?,status=?,deployed_by=?,deployed_at=? WHERE id=?`,e.Service,e.Version,e.Environment,e.Strategy,e.ActiveSlot,e.Status,e.DeployedBy,e.DeployedAt,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM deployments WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM deployments`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Deployment{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (1=0)"
        
    }
    if v,ok:=filters["environment"];ok&&v!=""{where+=" AND environment=?";args=append(args,v)}
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,service,version,environment,strategy,active_slot,status,deployed_by,deployed_at,created_at FROM deployments WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Deployment;for rows.Next(){var e Deployment;rows.Scan(&e.ID,&e.Service,&e.Version,&e.Environment,&e.Strategy,&e.ActiveSlot,&e.Status,&e.DeployedBy,&e.DeployedAt,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM deployments GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
