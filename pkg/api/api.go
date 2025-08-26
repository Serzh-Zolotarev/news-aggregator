package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"news-aggregator/pkg/storage"
	"strconv"
)

// Программный интерфейс сервера
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (a *API) endpoints() {
	a.router.Use(requestIdMiddleware)
	// получить новости
	a.router.HandleFunc("/news", a.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	// получить n новость
	a.router.HandleFunc("/news/{n}", a.postHandler).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	a.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
	a.router.Use(loggingMiddleware)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (a *API) Router() *mux.Router {
	return a.router
}

// Получение всех публикаций.
func (a *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	s := query.Get("s")
	page := query.Get("page")

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	posts, err := a.db.Posts(p, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// Получение публикации
func (a *API) postHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := a.db.Post(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// Добавление публикации.
func (a *API) addPostHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.NewsShortDetailed
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = a.db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обновление публикации.
func (a *API) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.NewsShortDetailed
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = a.db.UpdatePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление публикации.
func (a *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	var p storage.NewsShortDetailed
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = a.db.DeletePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
