package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-switchman/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Service{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var svc store.Service;json.NewDecoder(r.Body).Decode(&svc);if svc.Name==""{writeError(w,400,"name required");return};if svc.ActiveSlot==""{svc.ActiveSlot="blue"};if svc.SplitPercent==0{svc.SplitPercent=100};s.db.Create(&svc);writeJSON(w,201,svc)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleCutover(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{ToSlot string `json:"to_slot"`;Reason string `json:"reason"`};json.NewDecoder(r.Body).Decode(&req);if req.ToSlot==""{writeError(w,400,"to_slot required");return};s.db.Cutover(id,req.ToSlot,req.Reason);writeJSON(w,200,map[string]string{"status":"cutover complete"})}
func(s *Server)handleSplit(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Percent int `json:"percent"`};json.NewDecoder(r.Body).Decode(&req);s.db.SetSplit(id,req.Percent);writeJSON(w,200,map[string]string{"status":"split updated"})}
func(s *Server)handleHistory(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListCutovers(id);if list==nil{list=[]store.Cutover{}};writeJSON(w,200,list)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
