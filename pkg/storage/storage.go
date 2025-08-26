package storage

type NewsList struct {
	Pagination Pagination
	News       []NewsShortDetailed
}

type Pagination struct {
	Page       int
	MaxPages   int
	NewsOnPage int
}

// Публикация, получаемая из RSS.
type NewsShortDetailed struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts(n int, filter string) (*NewsList, error) // получение всех публикаций
	Post(n int) (*NewsShortDetailed, error)        // получение публикации
	AddPosts([]NewsShortDetailed) error            // создание новых публикаций
	AddPost(NewsShortDetailed) error               // создание новой публикации
	UpdatePost(NewsShortDetailed) error            // обновление публикации
	DeletePost(NewsShortDetailed) error            // удаление публикации по ID
}
