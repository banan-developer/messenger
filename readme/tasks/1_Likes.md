### 1. Создание таблицы лайков


```sql
CREATE TABLE `wall_likes` (
  `user_id` int NOT NULL,
  `post_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `post_id`),
  CONSTRAINT `fk_likes_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_likes_wall` FOREIGN KEY (`post_id`) REFERENCES `wall` (`idwall`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

```

### 2. Добавление и удаление лайка


**Поставить лайк:**

```sql
INSERT IGNORE INTO `wall_likes` (`user_id`, `post_id`) VALUES (21, 111);

```
**Убрать лайк:**

```sql
DELETE FROM `wall_likes` WHERE `user_id` = 21 AND `post_id` = 111;

```



**Лайки на посту:**

```sql
SELECT COUNT(*) as likes_count 
FROM `wall_likes` 
WHERE `post_id` = 111;

```

*Получить список всех постов вместе с количеством лайков (используя `LEFT JOIN`, чтобы вывелись даже посты без лайков):*

```sql
SELECT 
    w.idwall, 
    w.title, 
    w.text, 
    COUNT(l.user_id) as likes_count
FROM `wall` w
LEFT JOIN `wall_likes` l ON w.idwall = l.post_id
GROUP BY w.idwall;

```

### 4. Тестовы данные:

```SQL
-- 1. Добавление тестовых пользователей
-- Используем id от 901, чтобы не пересекаться с вашим текущим дампом
INSERT IGNORE INTO `users` (`id`, `login`, `password`, `name`, `about`, `avatar_url`, `sex`, `avatar_img`) VALUES
(901, 'testuser1@example.com', '$2a$10$dummyhash12345678901234', 'Иван Тестеров', 'Люблю тестировать функционал', 'unknown', 'Мужской', 'unknown'),
(902, 'testuser2@example.com', '$2a$10$dummyhash12345678901234', 'Анна Скриптова', 'Пишу SQL запросы и радуюсь', 'unknown', 'Женский', 'unknown'),
(903, 'testuser3@example.com', '$2a$10$dummyhash12345678901234', 'Гофер Гоферович', 'Разработчик на Go', 'unknown', 'Гофер', 'unknown');

-- 2. Добавление тестовых постов на стену
-- Привязываем их к созданным выше пользователям (users_id)
INSERT IGNORE INTO `wall` (`idwall`, `title`, `text`, `users_id`, `img_scr`) VALUES
(901, 'Первый пост Ивана', 'Привет, THE NOMAX! Это мой первый тестовый пост.', 901, NULL),
(902, 'Работа с БД', 'Ура, добавили таблицу лайков! Давайте тестировать.', 902, NULL),
(903, 'Про гоферов', 'Менеджмент памяти в Go просто великолепен.', 903, NULL),
(904, 'Второй пост Ивана', 'Продолжаю проверять, как отображаются записи на стене.', 901, NULL);

-- 3. Добавление тестовых лайков
-- Устанавливаем связи между пользователями и постами
INSERT IGNORE INTO `wall_likes` (`user_id`, `post_id`) VALUES
(902, 901), -- Анна лайкнула первый пост Ивана
(903, 901), -- Гофер лайкнул первый пост Ивана (у поста 901 теперь 2 лайка)
(901, 902), -- Иван лайкнул пост Анны
(903, 902), -- Гофер лайкнул пост Анны (у поста 902 теперь 2 лайка)
(902, 903), -- Анна лайкнула пост Гофера (у поста 903 теперь 1 лайк)
(901, 904); -- Иван лайкнул собственный пост (у поста 904 теперь 1 лайк)
```

### 5. Получить пость по id + колличество лайков

```sql
SELECT 
    w.idwall, 
    w.title, 
    w.text, 
    w.img_scr,
    w.created_at,
    COUNT(l.user_id) as likes_count
FROM `wall` w
LEFT JOIN `wall_likes` l ON w.idwall = l.post_id
WHERE w.idwall = 901  -- Замените 111 на нужный ID поста
GROUP BY 
    w.idwall, 
    w.title, 
    w.text,
    w.img_scr,
    w.created_at;
```

### 7. Как внедрить это в ваш код на Go (за адекватность не ручаюсь)

Чтобы это заработало в вашем бэкенде, вам нужно обновить структуру данных и добавить новый метод в репозиторий.

**Шаг 1. Обновление структуры в `internal/domain/Post.go**`
Добавьте поле `LikesCount` в структуру `WallPost`:

```go
package domain

type WallPost struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Img        string `json:"img"`
	CreatedAt  string `json:"created_at"`
	LikesCount int    `json:"likes_count"` // <--- Новое поле
}
```

**Шаг 2. Добавление метода в `internal/repository/Post.go**`
Напишите метод, который будет выполнять этот запрос и возвращать один пост по его ID:

```go
func (w *WallRepo) GetPostByID(PostID int) (*domain.WallPost, error) {
	query := `
		SELECT 
			w.idwall, 
			w.title, 
			w.text, 
			COALESCE(w.img_scr, ''), 
			COALESCE(w.created_at, ''),
			COUNT(l.user_id) as likes_count
		FROM wall w
		LEFT JOIN wall_likes l ON w.idwall = l.post_id
		WHERE w.idwall = ?
		GROUP BY 
			w.idwall, 
			w.title, 
			w.text,
			w.img_scr,
			w.created_at
	`

	var post domain.WallPost
	
	err := w.db.QueryRow(query, PostID).Scan(
		&post.Id, 
		&post.Title, 
		&post.Text, 
		&post.Img, 
		&post.CreatedAt,
		&post.LikesCount,
	)

	if err != nil {
		log.Printf("БД: Ошибка при получении поста по ID: %v", err)
		return nil, err
	}

	return &post, nil
}
```