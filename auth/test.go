package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	ory "github.com/ory/kratos-client-go"
)

type App struct {
	ory *ory.APIClient
}

func Auth() {
	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{
		{
			URL: fmt.Sprintf("http://localhost:4433"),
		},
	}

	app := &App{
		ory: ory.NewAPIClient(c),
	}
	mux := http.NewServeMux()

	// dashboard
	mux.Handle("/login", app.ensureCookieFlowID("login", app.handleLogin))
	mux.HandleFunc("/registration", app.ensureCookieFlowID("registration", app.handleRegister))
	mux.HandleFunc("/error", app.handleError)

	fmt.Println("Application launched and running on http://127.0.0.1:4455")
	// start the server
	http.ListenAndServe(":4455", mux)
}

func (s *App) handleRegister(w http.ResponseWriter, r *http.Request, cookie, flowID string) {
	// get the login flow
	flow, res, err := s.ory.FrontendApi.GetRegistrationFlow(r.Context()).Id(flowID).Cookie(cookie).Execute()
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusUnauthorized, err)
	}

	fmt.Println(flow.Ui)
	fmt.Println(res)

}

func (s *App) handleLogin(w http.ResponseWriter, r *http.Request, cookie, flowID string) {
	// get the login flow
	flow, _, err := s.ory.FrontendApi.GetLoginFlow(r.Context()).Id(flowID).Cookie(cookie).Execute()
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	b, err := json.Marshal(flow)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(b)
}

func (s *App) ensureCookieFlowID(flowType string, next func(w http.ResponseWriter, r *http.Request, cookie, flowID string)) http.HandlerFunc {
	// create redirect url based on flow type
	redirectURL := fmt.Sprintf("http://localhost:4433/self-service/%s/browser", flowType)

	return func(w http.ResponseWriter, r *http.Request) {
		// get flowID from url query parameters
		flowID := r.URL.Query().Get("flow")
		// if there is no flow id in url query parameters, create a new flow
		if flowID == "" {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		// get cookie from headers
		cookie := r.Header.Get("cookie")
		// if there is no cookie in header, return error
		if cookie == "" {
			writeError(w, http.StatusBadRequest, errors.New("missing cookie"))
			return
		}

		// call next handler
		next(w, r, cookie, flowID)
	}
}

func (s *App) handleError(w http.ResponseWriter, r *http.Request) {
	// get url query parameters
	errorID := r.URL.Query().Get("id")
	// get error details
	errorDetails, _, err := s.ory.FrontendApi.GetFlowError(r.Context()).Id(errorID).Execute()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	// marshal errorDetails to json
	errorDetailsJSON, err := json.MarshalIndent(errorDetails, "", "  ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(errorDetailsJSON)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_, e := w.Write([]byte(err.Error()))
	if e != nil {
		log.Fatal(err)
	}
}

//
//// save the cookies for any upstream calls to the Ory apis
//func withCookies(ctx context.Context, v string) context.Context {
//	return context.WithValue(ctx, "req.cookies", v)
//}
//
//func getCookies(ctx context.Context) string {
//	return ctx.Value("req.cookies").(string)
//}
//
//// save the session to display it on the dashboard
//func withSession(ctx context.Context, v *ory.Session) context.Context {
//	return context.WithValue(ctx, "req.session", v)
//}
//
//func getSession(ctx context.Context) *ory.Session {
//	return ctx.Value("req.session").(*ory.Session)
//}
//
//func (app *App) sessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		log.Printf("handling middleware request\n")
//
//		// set the cookies on the ory client
//		var cookies string
//
//		// this example passes all request.Cookies
//		// to `ToSession` function
//		//
//		// However, you can pass only the value of
//		// ory_session_projectid cookie to the endpoint
//		cookies = request.Header.Get("Cookie")
//
//		// check if we have a session
//		session, _, err := app.ory.FrontendApi.ToSession(request.Context()).Cookie(cookies).Execute()
//		if (err != nil && session == nil) || (err == nil && !*session.Active) {
//			// this will redirect the user to the managed Ory Login UI
//			http.Redirect(writer, request, "/.ory/self-service/login/browser", http.StatusSeeOther)
//			return
//		}
//
//		ctx := withCookies(request.Context(), cookies)
//		ctx = withSession(ctx, session)
//
//		// continue to the requested page (in our case the Dashboard)
//		next.ServeHTTP(writer, request.WithContext(ctx))
//		return
//	}
//}
//
//func (app *App) dashboardHandler() http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		tmpl, err := template.New("index.html").ParseFiles("index.html")
//		if err != nil {
//			http.Error(writer, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		session, err := json.Marshal(getSession(request.Context()))
//		if err != nil {
//			http.Error(writer, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		err = tmpl.ExecuteTemplate(writer, "index.html", string(session))
//		if err != nil {
//			http.Error(writer, err.Error(), http.StatusInternalServerError)
//			return
//		}
//	}
//}
