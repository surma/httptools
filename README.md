Package handlerlist provides the type `L` with which multiple
http.Handler can be chained to be executed sequentially.

# Example
```Go
func userData(w http.ResponseWriter, r *http.Request) {
	// Session magic
	session, err := openSession(r)
	if err != nil {
		http.Error(w, "Could not open session", http.StatusBadRequest)
	}
	w.(*handlerlist.VarsResponseWriter).Vars["UID"] = session.UserId
}

func showProfile(w http.ResponseWriter, r *http.Request) {
	uid := w.(*handlerlist.VarsResponseWriter).Vars["UID"].(string)
	profile := userProfile(uid)
	renderProfileTemplate(w, profile)
}

func main() {
	// ...
	http.Handle("/profile", handlerList.L {
		http.HandlerFunc(userData),
		handlerlist.SilentHandler(
			http.HandlerFunc(showProfile),
		)
	})
	// ...
}
```

