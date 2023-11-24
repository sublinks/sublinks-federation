package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"participating-online/sublinks-federation/internal/activitypub"
	"participating-online/sublinks-federation/internal/lemmy"

	"github.com/gorilla/mux"
	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/pressly/goose/v3"
)

//go:embed db/migration/*.sql
var embedMigrations embed.FS

func saveUser(user *activitypub.User, privateKey string) {
	dbUrl, _ := os.LookupEnv("DB_URL")
	db, err := getDb(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	res, err := db.Exec(fmt.Sprintf("INSERT INTO users (name, public_key, private_key) VALUES ('%s', '%s', '%s', '%s');", user.Name, user.Publickey.Keyid, user.Publickey.PublicKeyPem, privateKey))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)
}

func getDb(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("libsql", dbUrl)
	return db, err
}

func getLemmyClient(ctx context.Context) *lemmy.Client {
	user, _ := os.LookupEnv("LEMMY_USER")
	pass, _ := os.LookupEnv("LEMMY_PASSWORD")
	return lemmy.NewClient("https://demo.sublinks.org", user, pass)
}

func getUser(name string) (*activitypub.User, error) {
	dbUrl, _ := os.LookupEnv("DB_URL")
	db, err := getDb(dbUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("select * from users where name='%s';", name))
	if err != nil {
		return nil, err
	}
	var user activitypub.User
	if !rows.Next() {
		privatekey, publickey := GenerateKeyPair()
		user := activitypub.NewUser(name, privatekey, publickey)
		saveUser(&user, activitypub.GetPrivateKeyString(privatekey))
		return &user, nil
	}
	var id int
	var public_key, private_key string
	err = rows.Scan(&id, &name, &public_key, &private_key)
	if err != nil {
		fmt.Print(err)
		privatekey, publickey := GenerateKeyPair()
		user := activitypub.NewUser(name, privatekey, publickey)
		saveUser(&user, activitypub.GetPrivateKeyString(privatekey))
		return &user, nil
	}
	user = activitypub.NewExistingUser(name, private_key, public_key)
	return &user, nil
}

func runMigrations() {
	dbUrl, _ := os.LookupEnv("DB_URL")
	db, err := getDb(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "db/migration"); err != nil {
		panic(err)
	}
}

func main() {
	runMigrations()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Location", "/post/1")
		w.WriteHeader(http.StatusFound) // 302
	})
	r.HandleFunc("/users/{user}", GetUserInfoHandler).Methods("GET")
	r.HandleFunc("/users/{user}/inbox", GetInboxHandler).Methods("GET")
	r.HandleFunc("/users/{user}/inbox", PostInboxHandler).Methods("POST")
	r.HandleFunc("/users/{user}/outbox", GetOutboxHandler).Methods("GET")
	r.HandleFunc("/users/{user}/outbox", PostOutboxHandler).Methods("POST")
	r.HandleFunc("/post/{postId}", GetPostHandler).Methods("GET")

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func convertToApub(p *lemmy.Response) activitypub.Post {
	return activitypub.NewPost(
		p.PostView.Post.ApId,
		fmt.Sprintf("https://demo.sublinks.org/u/%s", p.PostView.Creator.Name),
		p.CommunityView.Community.ActorId,
		p.PostView.Post.Name,
		p.PostView.Post.Body,
		p.PostView.Post.Nsfw,
		p.PostView.Post.Published,
	)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := getLemmyClient(ctx)
	post, err := c.GetPost(ctx, vars["postId"])
	if err != nil {
		log.Println("Error reading post", err)
		return
	}
	postLd := convertToApub(post)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	w.Write(content)
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["user"] != "lazyguru" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User not found")
		return
	}

	user, err := getUser(vars["user"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(user, "", "  ")
	w.Write(content)
}

func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		os.Exit(1)
	}
	publickey := &privatekey.PublicKey
	return privatekey, publickey
}

func GetInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func PostInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func GetOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func PostOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}
