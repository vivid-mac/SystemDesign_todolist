package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	IsDone    bool      `db:"is_done"`
	Explanation string  `db:"explanation"`
	Priority  string    `db:"priority"`
	Deadline  time.Time `db:"deadline"`
	Tag       bool      `db:"tag"`
	Category  string    `db:"category"`
}

type User struct {
    ID        uint64    `db:"id"`
    Name      string    `db:"name"`
    Password  []byte    `db:"password"`
}

type Ownership struct {
	User_ID   uint64    `db:"user_id"`
	Task_ID   uint64    `db:"task_id"`
}

type Group struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type Belong struct {
	Group_ID  uint64    `db:"group_id"`  	
	User_ID   uint64    `db:"user_id"`
}

type Grouptask struct {
	Group_ID  uint64   `db:"group_id"`
	Task_ID   uint64    `db:"task_id"`
}