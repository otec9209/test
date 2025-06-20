package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib" // драйвер для PostgreSQL
)

// Функция для привязки меток к задачам
func linkLabelsToTasks(db *sql.DB) error {
	query := `
        INSERT INTO tasks_labels (task_id, label_id)
        SELECT t.id, l.id
        FROM tasks t
        CROSS JOIN labels l WHERE (
            (t.title = 'Помыть посуду' AND l.name = 'высокая срочность') OR
            (t.title = 'Вынести мусор' AND l.name = 'низкая срочность') OR
            (t.title = 'Купить продукты' AND l.name = 'средняя срочность') OR
            (t.title = 'Сделать уборку' AND l.name = 'высокая срочность') OR
            (t.title = 'Подготовить отчет' AND l.name = 'средняя срочность') OR
            (t.title = 'Отправить презентацию' AND l.name = 'средняя срочность') OR
            (t.title = 'Записаться на встречу' AND l.name = 'высокая срочность') OR
            (t.title = 'Созвон с командой' AND l.name = 'средняя срочность') OR
            (t.title = 'Проверить почту' AND l.name = 'низкая срочность') OR
            (t.title = 'Сделать бэкап' AND l.name = 'высокая срочность')
        )
        ON CONFLICT (task_id, label_id) DO NOTHING;
    `

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при привязке меток: %v", err)
	}

	log.Println("✅ Все задачи получили свои метки")
	return nil
}

// Очистка и перезалив данных
func clearAndReset(db *sql.DB) {
	_, _ = db.Exec("DELETE FROM tasks_labels")
	_, _ = db.Exec("DELETE FROM labels")
	_, _ = db.Exec("DELETE FROM tasks")
	_, _ = db.Exec("DELETE FROM usersStaff")

	fmt.Println("🧹 Все таблицы очищены")
}

// функция для сортировки задач по меткам
func getTasksByLabel(db *sql.DB, labelName string) {
	rows, err := db.Query(`
        SELECT 
            t.id,
            t.title,
            author.name AS author_name,
            assigned.name AS assigned_name
        FROM tasks t
        JOIN tasks_labels tl ON t.id = tl.task_id
        JOIN labels l ON tl.label_id = l.id
        JOIN usersStaff author ON t.author_id = author.id
        JOIN usersStaff assigned ON t.assigned_id = assigned.id
        WHERE l.name = $1
        ORDER BY t.id;
    `, labelName)

	if err != nil {
		log.Fatalf("Ошибка при выборке задач по метке '%s': %v\n", labelName, err)
	}
	defer rows.Close()

	fmt.Printf("\n📋 Задачи с меткой '%s':\n", labelName)

	for rows.Next() {
		var (
			id           int
			title        string
			authorName   string
			assignedName sql.NullString
		)
		err := rows.Scan(&id, &title, &authorName, &assignedName)
		if err != nil {
			log.Fatal("Ошибка при чтении данных:", err)
		}

		fmt.Printf("ID=%d | Название='%s' | Автор='%s' | Исполнитель='%s'\n",
			id, title, authorName, assignedName.String)
	}
}

// функция для поиска задач по именам
func getTasksByAuthor(db *sql.DB, authorName string) {
	rows, err := db.Query(`
        SELECT 
            t.id,
            t.title,
            author.name AS author_name,
            assigned.name AS assigned_name,
            l.name AS label_name
        FROM tasks t
        JOIN usersStaff author ON t.author_id = author.id
        LEFT JOIN usersStaff assigned ON t.assigned_id = assigned.id
        LEFT JOIN tasks_labels tl ON t.id = tl.task_id
        LEFT JOIN labels l ON tl.label_id = l.id
        WHERE author.name = $1
        ORDER BY t.id;
    `, authorName)

	if err != nil {
		log.Fatalf("Ошибка при выборке задач для автора '%s': %v\n", authorName, err)
	}
	defer rows.Close()

	fmt.Printf("\n📋 Задачи автора '%s':\n", authorName)

	for rows.Next() {
		var (
			id           int
			title        string
			authorName   string
			assignedName sql.NullString
			labelName    sql.NullString
		)
		err := rows.Scan(&id, &title, &authorName, &assignedName, &labelName)
		if err != nil {
			log.Fatal("Ошибка при чтении данных:", err)
		}

		fmt.Printf("ID=%d | Задача='%s' | Исполнитель='%s' | Метка='%s'\n",
			id, title, assignedName.String, labelName.String)
	}
}

func main() {
	// 1. Строка подключения
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=tasks sslmode=disable"

	// 2. Подключение к БД
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer db.Close()

	// 3. Проверка соединения
	err = db.Ping()
	if err != nil {
		log.Fatal("Нет связи с БД:", err)
	}

	fmt.Println("Подключено к БД!")

	clearAndReset(db)

	// 6. Таблица меток (labels)
	createTableSQL2 := `
    CREATE TABLE IF NOT EXISTS labels (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL2)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы 'labels':", err)
	}
	fmt.Println("✅ Таблица 'labels' создана или уже существует")

	// 8. Таблица сотрудников (usersStaff)
	createTableSQL3 := `
    CREATE TABLE IF NOT EXISTS usersStaff (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL3)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы 'usersStaff':", err)
	}
	fmt.Println("✅ Таблица 'usersStaff' создана или уже существует")

	// 10. Таблица задач (tasks)
	createTableSQL4 := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT,
        opened BIGINT NOT NULL DEFAULT extract(epoch FROM now()),
        closed BIGINT,
        author_id INTEGER REFERENCES usersStaff(id) NOT NULL,
        assigned_id INTEGER REFERENCES usersStaff(id)
    );`

	_, err = db.Exec(createTableSQL4)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы 'tasks':", err)
	}
	fmt.Println("✅ Таблица 'tasks' создана или уже существует")

	// 12. Связующая таблица (tasks_labels)
	createTableSQL5 := `
    CREATE TABLE IF NOT EXISTS tasks_labels (
        task_id INTEGER REFERENCES tasks(id),
        label_id INTEGER REFERENCES labels(id),
        PRIMARY KEY (task_id, label_id)
    );`

	_, err = db.Exec(createTableSQL5)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы 'tasks_labels':", err)
	}
	fmt.Println("✅ Таблица 'tasks_labels' создана или уже существует")

	var countLabels int
	err = db.QueryRow("SELECT COUNT(*) FROM labels").Scan(&countLabels)
	if err != nil {
		log.Fatal("Ошибка при проверке количества меток:", err)
	}

	if countLabels == 0 {
		createLabelsSQL := `
        INSERT INTO labels (name) VALUES 
            ('низкая срочность'), 
            ('средняя срочность'), 
            ('высокая срочность');
    `

		_, err = db.Exec(createLabelsSQL)
		if err != nil {
			log.Fatal("Ошибка при добавлении меток:", err)
		}
		fmt.Println("✅ Метки успешно добавлены")
	} else {
		fmt.Printf("В базе уже есть %d меток, новые метки не добавляются\n", countLabels)
	}

	//Добавить сотрудников в таблицу usersStaff
	var countStaff int
	err = db.QueryRow("SELECT COUNT(*) FROM usersStaff").Scan(&countStaff)
	if err != nil {
		log.Fatal("Ошибка при проверке количества сотрудников:", err)
	}

	if countStaff == 0 {
		createUserStaffSQL := `
	INSERT INTO usersStaff (name)
	VALUES ('Алексей'), ('Иван'), ('Мария'), ('Анна'), ('Дмитрий');`

		_, err = db.Exec(createUserStaffSQL)
		if err != nil {
			log.Fatal("Ошибка при добавлении сотрудников:", err)
		}
		fmt.Println("✅ 5 сотрудников успешно добавлены")
	} else {
		fmt.Printf("В базе уже есть %d сотрудников, новые сотрудники не добавляются\n", countStaff)
	}

	// Проверяем, есть ли уже задачи в базе
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		log.Fatal("Ошибка при проверке количества задач:", err)
	}

	if count == 0 {
		// 13. создание новых задач только если их нет

		createTaskSQL := `
INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Помыть посуду', 'Нужно помыть кружки и тарелки', 
       (SELECT id FROM usersStaff WHERE name = 'Алексей'), 
       (SELECT id FROM usersStaff WHERE name = 'Иван')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Помыть посуду');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Вынести мусор', 'Вывести мусор из кухни и ванной',
       (SELECT id FROM usersStaff WHERE name = 'Алексей'), 
       (SELECT id FROM usersStaff WHERE name = 'Иван')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Вынести мусор');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Купить продукты', 'Молоко, хлеб, яйца',
       (SELECT id FROM usersStaff WHERE name = 'Алексей'), 
       (SELECT id FROM usersStaff WHERE name = 'Иван')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Купить продукты');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Сделать уборку', 'Протереть пыль и пропылесосить',
       (SELECT id FROM usersStaff WHERE name = 'Иван'), 
       (SELECT id FROM usersStaff WHERE name = 'Алексей')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Сделать уборку');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Подготовить отчет', 'За прошлый месяц',
       (SELECT id FROM usersStaff WHERE name = 'Иван'), 
       (SELECT id FROM usersStaff WHERE name = 'Мария')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Подготовить отчет');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Отправить презентацию', 'Клиенту до обеда',
       (SELECT id FROM usersStaff WHERE name = 'Мария'), 
       (SELECT id FROM usersStaff WHERE name = 'Алексей')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Отправить презентацию');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Записаться на встречу', 'С менеджером проекта',
       (SELECT id FROM usersStaff WHERE name = 'Мария'), 
       (SELECT id FROM usersStaff WHERE name = 'Иван')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Записаться на встречу');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Созвон с командой', 'Обсудить план недели',
       (SELECT id FROM usersStaff WHERE name = 'Алексей'), 
       (SELECT id FROM usersStaff WHERE name = 'Мария')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Созвон с командой');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Проверить почту', 'Все непрочитанные письма',
       (SELECT id FROM usersStaff WHERE name = 'Иван'), 
       (SELECT id FROM usersStaff WHERE name = 'Мария')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Проверить почту');

INSERT INTO tasks (title, content, author_id, assigned_id)
SELECT 'Сделать бэкап', 'Бэкап всех данных',
       (SELECT id FROM usersStaff WHERE name = 'Мария'), 
       (SELECT id FROM usersStaff WHERE name = 'Анна')
WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE title = 'Сделать бэкап');
`

		_, err = db.Exec(createTaskSQL)
		if err != nil {
			log.Fatal("Ошибка при добавлении задач:", err)
		}
		fmt.Println("✅ 10 задач успешно добавлены")
	} else {
		fmt.Printf("В базе уже есть %d задач, новые задачи не добавляются\n", count)
	}

	// 14. Получить список всех задач с именами сотрудников
	getTasksSQL := `
SELECT 
    t.id,
    t.title,
    t.content,
    t.opened,
    t.closed,
    author.name AS author_name,
    assigned.name AS assigned_name,
    l.name AS label_name
FROM tasks t
LEFT JOIN usersStaff author ON t.author_id = author.id
LEFT JOIN usersStaff assigned ON t.assigned_id = assigned.id
LEFT JOIN tasks_labels tl ON t.id = tl.task_id
LEFT JOIN labels l ON tl.label_id = l.id
ORDER BY t.id;
`

	err = linkLabelsToTasks(db)
	if err != nil {
		log.Println(err)
	}

	// Потом читаем всё вместе
	fmt.Println("\nСписок всех задач:")
	rows, err := db.Query(getTasksSQL)
	if err != nil {
		log.Fatal("Ошибка при получении списка задач:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            int
			title         string
			content       sql.NullString
			opened        sql.NullInt64
			closed        sql.NullInt64
			author_name   string
			assigned_name sql.NullString
			label_name    sql.NullString
		)
		err = rows.Scan(&id, &title, &content, &opened, &closed, &author_name, &assigned_name, &label_name)
		if err != nil {
			log.Fatal("Ошибка при чтении данных задачи:", err)
		}
		fmt.Printf("ID=%d | Задача: %s | Автор: %s | Исполнитель: %s | Метка: %s\n",
			id, title, author_name, assigned_name.String, label_name.String)
	}

	// Выберите задачи по метке
	getTasksByLabel(db, "высокая срочность")
	getTasksByLabel(db, "низкая срочность")
	getTasksByLabel(db, "средняя срочность")

	// Вывести задачи только от Алексея
	getTasksByAuthor(db, "Алексей")
	// Вывести задачи только от Ивана
	getTasksByAuthor(db, "Иван")

}
