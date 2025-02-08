package handlers

import "github.com/gorilla/sessions"

// Store is the session store used for authentication
var Store = sessions.NewCookieStore([]byte("your-secret-key"))
