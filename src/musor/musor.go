package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib" // драйвер для PostgreSQL
)

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

	// 4. Таблица пользователей (users)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы 'users':", err)
	}
	fmt.Println("✅ Таблица 'users' создана или уже существует")

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
		VALUES 
		('Помыть посуду', 'Нужно помыть кружки и тарелки', 1, 2),
		('Вынести мусор', 'Вывести мусор из кухни и ванной', 1, 2),
		('Купить продукты', 'Молоко, хлеб, яйца', 1, 2),
		('Сделать уборку', 'Протереть пыль и пропылесосить', 2, 1),
		('Подготовить отчет', 'За прошлый месяц', 2, 3),
		('Отправить презентацию', 'Клиенту до обеда', 3, 1),
		('Записаться на встречу', 'С менеджером проекта', 3, 2),
		('Созвон с командой', 'Обсудить план недели', 1, 3),
		('Проверить почту', 'Все непрочитанные письма', 2, 3),
		('Сделать бэкап', 'Бэкап всех данных', 3, 1);`

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
		author.name as author_name,
		assigned.name as assigned_name
	FROM tasks t
	LEFT JOIN usersStaff author ON t.author_id = author.id
	LEFT JOIN usersStaff assigned ON t.assigned_id = assigned.id
	ORDER BY t.id;`

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
		)
		err = rows.Scan(&id, &title, &content, &opened, &closed, &author_name, &assigned_name)
		if err != nil {
			log.Fatal("Ошибка при чтении данных задачи:", err)
		}
		fmt.Printf("ID=%d, Название='%s', Автор='%s', Исполнитель='%s'\n",
			id, title, author_name, assigned_name.String)
	}
}
