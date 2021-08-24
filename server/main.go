package main

import (
	"SessionCookies/token"
	"SessionCookies/trie"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	//"github.com/stripe/stripe-go"
)

// stripe.Key = "pk_test_51IzvABIch0MObKLPR2lMcv5edvxKNpo4QCPA6JGtPaXAtNF8xQLHnkg1vGOeAXgRRBrfCZK0ODrTlfLxM9YcqfzD00JnYcW2H9"
// ch, err := charge.New(params)

// Declaration of global variables
var (
	tpl   *template.Template
	mutex sync.Mutex // Concurrency

	baseURL   = "https://localhost"
	httpsPort = ":5221"

	db *sql.DB

	// searches for
	maskNameSearch     *trie.Trie // mask names
	maskCategorySearch *trie.Trie // mask categories
	sellerSearch       *trie.Trie // seller username
)

func sizesBitAnd(a, b uint8) int {
	return int(a & b)
}

// for sellerEditProduct.gohtml
// note: only max 8 images are allowed to be uploaded for a mask
func generateEmptyImg(n int) []int {
	numberOfNewImgNeeded := 8 - n
	newArray := make([]int, numberOfNewImgNeeded)
	for i := 0; i < len(newArray); i++ {
		newArray[i] = n + i
	}
	return newArray
}

// see if that field is available on that struct passed to the template
// e.g. does Data struct contain a field called "Member"?
func avail(field string, data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(field).IsValid()
}

// makeLink replaces spaces with plus
func lowercase(field string) string {
	return strings.ToLower(field)
}

func init() {
	tpl = template.Must(template.New("").Funcs(template.FuncMap{
		"sizesbitand":      sizesBitAnd, // checks a&b
		"generateemptyimg": generateEmptyImg,
		"avail":            avail,
		"lowercase":        lowercase,
	}).ParseGlob("templates/*")) // Must() reads the templates

	// set time location to Singapore
	loc, err := time.LoadLocation("Asia/Singapore")
	// handle err
	if err != nil {
		log.Fatalln(err)
	}
	time.Local = loc // sets the global timezone
}

// This function initialise sessions from the saved csv file to the dbSession map which is a map of key cookie value, value username
func initBuyerSession() { // Convert CSV file 2D array to patient account map
	records := buyerSessionReadCsv()
	for i := 0; i < len(records); i++ {
		for j := 0; j < 2; j++ {
			cookieValue := records[i][0]
			un := records[i][1]
			dbBuyerSessions[cookieValue] = un
		}
	}
}

// This function initialise sessions from the saved csv file to the dbSession map which is a map of key cookie value, value username
func initSellerSession() { // Convert CSV file 2D array to patient account map
	records := sellerSessionReadCsv()
	for i := 0; i < len(records); i++ {
		for j := 0; j < 2; j++ {
			cookieValue := records[i][0]
			un := records[i][1]
			dbSellerSessions[cookieValue] = un
		}
	}
}

func home(res http.ResponseWriter, req *http.Request) {
	// // check for seller cookie; redirrect to /seller if found
	// _, err := askForSellerCookie(res, req)
	// if err == nil {
	// 	http.Redirect(res, req, "/seller", http.StatusFound)
	// }

	tpl.ExecuteTemplate(res, "homePage.gohtml", nil)
}

func buyerLoginSuccess(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "buyerLoginSuccess.gohtml", nil)
}

func sellerLoginSuccess(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "sellerLoginSuccess.gohtml", nil)
}

func paymentCheckout(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "paymentCheckout.html", nil)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		initBuyerSession()
	}()
	go func() {
		defer wg.Done()
		initSellerSession()
	}()
	go func() {
		defer wg.Done()
		// Generate key and convert to string for password reset feature
		key := generateSecretKey()
		keyStr = string(key)
		var err error
		maker, err = token.NewJWTMaker(keyStr)
		if err != nil {
			fmt.Println("Error generating token maker!")
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3307)/crazyformasks_db?parseTime=true")
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("crazyformasks_db opened.")
		}

		err = PrepareSearchTries(db)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Search trie created and loaded.")
		}
	}()
	wg.Wait()
	defer func() {
		db.Close()
		Info.Println("Connection to crazyformasks_db closed.")
	}()

	// TEST CODE
	// _, found := (search.Find("Alfian", 7))
	// fmt.Println(found)

	// ch2 := make(chan *trie.EntryInfo)
	// go search.TrieTraverserReturningNodes(search.Root, "", 7, ch2)
	// for ei := range ch2 {
	// 	fmt.Printf("%32s is marked %3d times as mask name, %3d times as mask category, %3d times as seller's username\n", ei.Entry, ei.Node.TimesAsName, ei.Node.TimesAsCategory, ei.Node.TimesAsSeller)
	// }
	// END TEST CODE

	// TODO: Cron job
	// Cron job 1: @midnight, calculate, for each mask, number sold last 30 days.

	// http multiplexer with gorilla/mux
	r := mux.NewRouter()

	r.HandleFunc("/", home) // signup login

	// load images
	images := http.StripPrefix("/img/", http.FileServer(http.Dir("./img/")))
	r.PathPrefix("/img/").Handler(images)
	css := http.StripPrefix("/css/", http.FileServer(http.Dir("./css/")))
	r.PathPrefix("/css/").Handler(css)
	js := http.StripPrefix("/js/", http.FileServer(http.Dir("./js/")))
	r.PathPrefix("/js/").Handler(js)

	r.HandleFunc("/adminlogout", adminLogout)
	r.HandleFunc("/buyerlogout", buyerLogout)
	r.HandleFunc("/sellerlogout", sellerLogout)
	r.HandleFunc("/buyersignup", buyerSignup)
	r.HandleFunc("/sellersignup", sellerSignup)
	r.HandleFunc("/adminlogin", adminLogin)
	r.HandleFunc("/buyerlogin", buyerLogin)
	r.HandleFunc("/sellerlogin", sellerLogin)
	r.HandleFunc("/adminresetpassword", adminResetPassword)
	r.HandleFunc("/adminresetchangepassword", adminResetChangePassword)
	r.HandleFunc("/buyerresetpassword", buyerResetPassword)
	r.HandleFunc("/buyerresetchangepassword", buyerResetChangePassword)
	r.HandleFunc("/buyertoken/{token}", resetBuyerPasswordLinkClicked)
	r.HandleFunc("/sellerresetpassword", sellerResetPassword)
	r.HandleFunc("/sellerresetchangepassword", sellerResetChangePassword)
	r.HandleFunc("/sellertoken/{token}", resetSellerPasswordLinkClicked)
	r.HandleFunc("/buyerloginsuccess", buyerLoginSuccess)
	r.HandleFunc("/sellerloginsuccess", sellerLoginSuccess)
	r.HandleFunc("/deletebuyer", deleteBuyer)
	r.HandleFunc("/deleteseller", deleteSeller)
	r.HandleFunc("/adminchangepassword", adminChangePassword)
	r.HandleFunc("/buyerchangepassword", buyerChangePassword)
	r.HandleFunc("/sellerchangepassword", sellerChangePassword)
	r.HandleFunc("/paymentcheckout", paymentCheckout)
	r.HandleFunc("/deletebuyersessions", deleteBuyerSessions)
	r.HandleFunc("/deletesellersessions", deleteSellerSessions)

	// autocomplete for search bar
	r.HandleFunc("/autocomplete", AutocompleteHandler).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/search", SearchHandler).
		Methods("GET")

	// handlers for sellers
	// assumes seller has logged in (check cookie in each handler)
	// /seller also shows seller products being sold
	seller := r.PathPrefix("/seller").Subrouter()
	// seller main menu (landing page for sellers)
	seller.HandleFunc("", SellerPortalHandler)
	seller.HandleFunc("/addproduct", TryAddProductHandler).
		Methods("GET")
	seller.HandleFunc("/addproduct", AddProductHandler).
		Methods("POST")
	seller.HandleFunc("/viewproduct/{maskid:[0-9]+}", ViewProductHandler).
		Methods("GET")
	seller.HandleFunc("/editproduct/{maskid:[0-9]+}", TryEditProductHandler).
		Methods("GET")
	seller.HandleFunc("/editproduct/{maskid:[0-9]+}", EditProductHandler).
		Methods("POST")
	seller.HandleFunc("/deleteproduct", DeleteProductHandler).
		Methods("POST")
	seller.HandleFunc("/autocomplete", AutocompleteHandler).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")
	seller.HandleFunc("/search", SearchHandler).
		Methods("GET")

	// Admin REST API functions for focus words
	r.HandleFunc("/admin/api/v1/login", LoginAndGetAPIKeyHandler).
		Methods("POST").
		HeadersRegexp("Authorization", "^Basic [A-Za-z0-9+/=]+$")
	r.HandleFunc("/admin/api/v1/resetapikey", RegenerateAPIKeyHandler).
		Methods("POST").
		HeadersRegexp("Authorization", "^Basic [A-Za-z0-9+/=]+$")

	r.HandleFunc("/admin/api/v1/focuswords", AllFocusWordsFromATrieHandler).
		Methods("GET").
		Queries("searchtrie", "{searchtrie:[1-3]}")
		// Queries("key", "{key:^[a-zA-Z0-9_-]{32}$}")

	admin := r.PathPrefix("/admin/api/v1/focusword").
		Subrouter()
	admin.HandleFunc("", AddFocusWordHandler).
		Methods("POST").
		Headers("Content-Type", "application/json")
	admin.HandleFunc("", InfoAboutFocusWordHandler).
		Methods("GET").
		Queries("searchtrie", "{searchtrie:[1-3]}" /*, "entry", "{entry:^[a-zA-Z0-9- ]{,32}$}"*/)
	admin.HandleFunc("", UpdateFocusWordHandler).
		Methods("PUT").
		Headers("Content-Type", "application/json")
	admin.HandleFunc("", NullifyFocusWordHandler).
		Methods("DELETE").
		Queries("searchtrie", "{searchtrie:[1-3]}" /*, "entry", "{entry:^[a-zA-Z0-9- ]{,32}$}"*/)

	// handlers for browsers/buyers
	r.HandleFunc("/browser", browserPage)
	r.HandleFunc("/product", productPage)
	r.HandleFunc("/review", reviewPage)
	r.HandleFunc("/review/add", reviewAdd)

	r.HandleFunc("/shoppingcart", shoppingCartPage)
	r.HandleFunc("/shoppingcart/add", shoppingCartAdd)
	r.HandleFunc("/shoppingcart/update", shoppingCartUpdate)
	r.HandleFunc("/shoppingcart/delete", shoppingCartDelete)

	r.HandleFunc("/checkout", checkoutPage)
	r.HandleFunc("/payment", paymentPage)
	r.HandleFunc("/orders", ordersPage)

	//added by Kum Choon @ 17-Jun-2021
	r.HandleFunc("/buyerinfochange", buyerInfoChange)
	r.HandleFunc("/sellerinfochange", sellerInfoChange)
	//end of added by Kum Choon @ 17-Jun-2021

	// handles redirection to HTTPS (if user tries to use HTTP)
	// redirectTLS defined in muxGeneral.go
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("The redirect to HTTPS has encountered a ListenAndServe error: %v", err)
		}
	}()

	// HTTPS Server
	go func() {
		/*if err := http.ListenAndServeTLS(httpsPort, "cert.pem", "key.pem", r); err != nil {
			log.Fatalf("Our HTTPS server encountered a ListenAndServe error: %v", err)
		}*/
		// from https://github.com/denji/golang-tls
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		srv := &http.Server{
			Addr:         httpsPort,
			Handler:      r,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
	}()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt) // ctrl-C to quit app
	// Block until a signal is received.
	<-c
	Info.Println("Program has shut down.")
}

/*
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kim3z/golang-rest-auth/controllers"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome home!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Home)
	router.POST("/api/user/create", controllers.CreateUser)
	router.POST("/api/user/login", controllers.LoginUser)
	router.POST("/api/user/forgot-password/:email", controllers.ForgotPassword)
	router.GET("/api/user/reset-psw-check/:reset-token", controllers.ResetPasswordCheck)
	router.POST("/api/user/reset-password", controllers.ResetPassword)

	// Protected route example
	// router.GET("/foo", middleware.Authenticate(httprouter.Handle(Foo)))

	fmt.Println("Listening on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
*/

// https://github.com/gorilla/mux  <--------------------------URL paths to handlers
