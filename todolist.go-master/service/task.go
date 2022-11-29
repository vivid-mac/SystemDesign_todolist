package service

import (
	"net/http"
	"strconv"
	"fmt"
    "time"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
	database "todolist.go/db"
)

//２つのTask構造体スライスの和集合を返す
func Union(task1 []database.Task, task2 []database.Task) []database.Task{
    Tasks := task1
    for i:=0; i<len(task2); i++{
        check := true
        for j:=0; j<len(task1); j++{
            if task2[i].ID == task1[j].ID {
                check = false //共通した要素ならfalseになる
                break
            }
        }
        if check { //共通していないなら加える
            Tasks = append(Tasks, task2[i])
        }
    }

    return Tasks
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
    userID := sessions.Default(ctx).Get("user")
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

    //Get query parameter
    kw := ctx.Query("kw")
	check_done_str := ctx.Query("check_done")
    check_done := false
    if check_done_str == "done"{
        check_done = true
    }
    category := ctx.Query("category")
    tag_str := ctx.Query("tag")
    tag := false
    if tag_str == "tagged"{
        tag = true
    }

	// Get tasks in DB
	var tasks []database.Task
    query := "SELECT id, title, created_at, is_done, priority, deadline, tag, category FROM tasks INNER JOIN ownership ON task_id = id WHERE user_id = ?"
	switch {
    case kw != "" && check_done_str != "" && category != "" && tag_str != "": //case1
        err = db.Select(&tasks, query + " AND title LIKE ? AND is_done = ? AND category = ? AND tag = ?", userID, "%" + kw + "%", check_done, category, tag)
    case kw != "" && check_done_str != "" && category != "" && tag_str == "": //case1.1
        err = db.Select(&tasks, query + " AND title LIKE ? AND is_done = ? AND category = ?", userID, "%" + kw + "%", check_done, category)
    case kw != "" && check_done_str != "" && category == "" && tag_str != "": //case1.2
        err = db.Select(&tasks, query + " AND title LIKE ? AND is_done = ? AND tag = ?", userID, "%" + kw + "%", check_done, tag)
    case kw != "" && check_done_str != "" && category == "" && tag_str == "": //case1.3
        err = db.Select(&tasks, query + " AND title LIKE ? AND is_done = ?", userID, "%" + kw + "%", check_done)      
    case kw != "" && check_done_str == "" && category != "" && tag_str != "": //case2
        err = db.Select(&tasks, query + " AND title LIKE ? AND category = ? AND tag = ?", userID, "%" + kw + "%", category, tag)
    case kw != "" && check_done_str == "" && category != "" && tag_str == "": //case2.1
        err = db.Select(&tasks, query + " AND title LIKE ? AND category = ? AND tag = ?", userID, "%" + kw + "%", category)
    case kw != "" && check_done_str == "" && category == "" && tag_str != "": //case2.2
        err = db.Select(&tasks, query + " AND title LIKE ? AND tag = ?", userID, "%" + kw + "%", category, tag)
    case kw != "" && check_done_str == "" && category == "" && tag_str == "": //case2.3
        err = db.Select(&tasks, query + " AND title LIKE ?", userID, "%" + kw + "%", category)      
    case kw == "" && check_done_str != "" && category != "" && tag_str != "": //case3
        err = db.Select(&tasks, query + " AND is_done = ? AND category = ? AND tag = ?", userID, check_done, category, tag)
    case kw == "" && check_done_str != "" && category != "" && tag_str == "": //case3.1
        err = db.Select(&tasks, query + " AND is_done = ? AND category = ?", userID, check_done, category)
    case kw == "" && check_done_str != "" && category == "" && tag_str != "": //case3.2
        err = db.Select(&tasks, query + " AND is_done = ? AND tag = ?", userID, check_done, tag)
    case kw == "" && check_done_str != "" && category == "" && tag_str == "": //case3.3
        err = db.Select(&tasks, query + " AND is_done = ?", userID, check_done)
    case kw == "" && check_done_str == "" && category != "" && tag_str != "": //case4
        err = db.Select(&tasks, query + " AND category = ? AND tag = ?", userID, category, tag)
    case kw == "" && check_done_str == "" && category != "" && tag_str == "": //case4.1
        err = db.Select(&tasks, query + " AND category = ?", userID, category)
    case kw == "" && check_done_str == "" && category == "" && tag_str != "": //case5
        err = db.Select(&tasks, query + " AND tag = ?", userID, tag)
    default: //case6
        err = db.Select(&tasks, query, userID)
    }
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

    //　ユーザがグループに所属しているか確認する
    var check int
    err = db.Get(&check, "SELECT COUNT(*) FROM belong WHERE user_id = ?", userID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    if check > 0 { //check=0のときグループに所属していない
        //もし所属しているならtasksスライスに共有しているタスクを加える
        var groupid uint64
        err = db.Get(&groupid, "SELECT group_id FROM belong WHERE user_id = ?", userID)
        if err != nil {
            Error(http.StatusBadRequest, err.Error())(ctx)
            return
        }
        var grouptasks []database.Task
        query = "SELECT id, title, created_at, is_done, priority, deadline, tag, category FROM tasks INNER JOIN share ON task_id = id WHERE group_id = ?"
        switch {
        case kw != "" && check_done_str != "" && category != "" && tag_str != "": //case1
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND is_done = ? AND category = ? AND tag = ?", groupid, "%" + kw + "%", check_done, category, tag)
        case kw != "" && check_done_str != "" && category != "" && tag_str == "": //case1.1
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND is_done = ? AND category = ?", groupid, "%" + kw + "%", check_done, category)
        case kw != "" && check_done_str != "" && category == "" && tag_str != "": //case1.2
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND is_done = ? AND tag = ?", groupid, "%" + kw + "%", check_done, tag)
        case kw != "" && check_done_str != "" && category == "" && tag_str == "": //case1.3
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND is_done = ?", groupid, "%" + kw + "%", check_done)      
        case kw != "" && check_done_str == "" && category != "" && tag_str != "": //case2
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND category = ? AND tag = ?", groupid, "%" + kw + "%", category, tag)
        case kw != "" && check_done_str == "" && category != "" && tag_str == "": //case2.1
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND category = ? AND tag = ?", groupid, "%" + kw + "%", category)
        case kw != "" && check_done_str == "" && category == "" && tag_str != "": //case2.2
            err = db.Select(&grouptasks, query + " AND title LIKE ? AND tag = ?", groupid, "%" + kw + "%", category, tag)
        case kw != "" && check_done_str == "" && category == "" && tag_str == "": //case2.3
            err = db.Select(&grouptasks, query + " AND title LIKE ?", groupid, "%" + kw + "%", category)      
        case kw == "" && check_done_str != "" && category != "" && tag_str != "": //case3
            err = db.Select(&grouptasks, query + " AND is_done = ? AND category = ? AND tag = ?", groupid, check_done, category, tag)
        case kw == "" && check_done_str != "" && category != "" && tag_str == "": //case3.1
            err = db.Select(&grouptasks, query + " AND is_done = ? AND category = ?", groupid, check_done, category)
        case kw == "" && check_done_str != "" && category == "" && tag_str != "": //case3.2
            err = db.Select(&grouptasks, query + " AND is_done = ? AND tag = ?", groupid, check_done, tag)
        case kw == "" && check_done_str != "" && category == "" && tag_str == "": //case3.3
            err = db.Select(&grouptasks, query + " AND is_done = ?", groupid, check_done)
        case kw == "" && check_done_str == "" && category != "" && tag_str != "": //case4
            err = db.Select(&grouptasks, query + " AND category = ? AND tag = ?", groupid, category, tag)
        case kw == "" && check_done_str == "" && category != "" && tag_str == "": //case4.1
            err = db.Select(&grouptasks, query + " AND category = ?", groupid, category)
        case kw == "" && check_done_str == "" && category == "" && tag_str != "": //case5
            err = db.Select(&grouptasks, query + " AND tag = ?", groupid, tag)
        default: //case6
            err = db.Select(&grouptasks, query, groupid)
        }
        if err != nil {
            Error(http.StatusInternalServerError, err.Error())(ctx)
            return
        }

        tasks = Union(tasks, grouptasks) //二つのスライスを結合する
    }

    // Get number of tasks
    taskNum := len(tasks)
    // Get number of page
    pageNum, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

    var tasks_split []database.Task
    if pageNum * 10 >= taskNum {
        tasks_split = tasks[pageNum*10-10:taskNum]
    } else {
        tasks_split = tasks[pageNum*10-10:pageNum*10]
    }
    page := fmt.Sprintf("%d", pageNum) //現在のページ
    prev := fmt.Sprintf("%d", pageNum-1) //一つ前のページ
    next := fmt.Sprintf("%d", pageNum+1) //次のページ
    if pageNum == 1 { //一つ目のページの時はprevはなし
        prev = "Nothing"
    }
    if pageNum * 10 >= taskNum { //最後のページの時はnextはなし
        next = "Nothing"
    }

    ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Page": page, "Prev": prev, "Next": next, "Title": "Task list", "Tasks": tasks_split, "Kw": kw, "Check_done": check_done, "Check_done_str": check_done_str, "Category": category, "Tag_str": tag_str})

	// Render tasks
	//ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Page": "1", "Title": "Task list", "Tasks": tasks, "Kw": kw, "Check_done": check_done, "Check_done_str": check_done_str, "Category": category, "Tag_str": tag_str})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

    // Get share information
    var group_id uint64
    err = db.Get(&group_id, "SELECT group_id FROM share WHERE task_id = ?", id)
    share := true
    if group_id == 0 { //group_id=0のときは共有しない
        share = false
    }

    // Calculate countdown
    countdown := int(task.Deadline.Sub(time.Now()).Hours() / 24)

	// Render task
	ctx.HTML(http.StatusOK, "task.html", gin.H{"Task": task, "Countdown": countdown, "Share": share})
}

func NewTaskForm(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "form_new_task.html", gin.H{"Title": "Task registration"})
}

func RegisterTask(ctx *gin.Context) {
    userID := sessions.Default(ctx).Get("user")
    // Get task title
    title, exist := ctx.GetPostForm("title")
    if !exist {
        Error(http.StatusBadRequest, "No title is given")(ctx)
        return
    }
    // Get task desctiption
    description, _ := ctx.GetPostForm("description")
    // Get task priority
    priority, exist := ctx.GetPostForm("priority")
    if !exist {
        Error(http.StatusBadRequest, "No priority is given")(ctx)
        return
    }
    // Get task deadline
    deadline, exist := ctx.GetPostForm("deadline")
    if !exist {
        Error(http.StatusBadRequest, "No deadline is given")(ctx)
        return
    }
    // Get task category
    category, exist := ctx.GetPostForm("category")
    if !exist {
        Error(http.StatusBadRequest, "No category is given")(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	// Create new data with given title on DB
    tx := db.MustBegin()
    result, err := tx.Exec("INSERT INTO tasks (title, explanation, priority, deadline, category) VALUES (?, ?, ?, ?, ?)", title, description, priority, deadline, category)
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    taskID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    _, err = tx.Exec("INSERT INTO ownership (user_id, task_id) VALUES (?, ?)", userID, taskID)
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // shareテーブルにデフォルトグループid＝０として登録する.(共有していない時はグループid=0)
    _, err = db.Exec("INSERT INTO share (group_id, task_id) VALUES (?, ?)", 0, taskID)
    if err != nil {
        tx.Rollback()
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	tx.Commit()
    ctx.Redirect(http.StatusFound, fmt.Sprintf("/task/%d", taskID))
}

func EditTaskForm(ctx *gin.Context) {
    // ID の取得
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Get target task
    var task database.Task
    err = db.Get(&task, "SELECT id, title, created_at, is_done, explanation, priority, deadline, tag, category FROM tasks WHERE id=?", id)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    //　ユーザがグループに所属しているか確認する
    userID := sessions.Default(ctx).Get("user")
    var check int
    err = db.Get(&check, "SELECT COUNT(*) FROM belong WHERE user_id = ?", userID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    belonged := false
    if check > 0 {//check>0ならばグループに所属している
        belonged = true
    }
    // Get task share information
    share := true
    var groupid int
    err = db.Get(&groupid, "SELECT group_id FROM share WHERE task_id = ?", task.ID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    if groupid == 0 {//もしgroupid=0なら共有していない
        share = false
    }

    // Render edit form
    ctx.HTML(http.StatusOK, "form_edit_task.html",
        gin.H{"Title": fmt.Sprintf("Edit task %d", task.ID), "Task": task, "Belonged": belonged, "Share": share})
}

func UpdateTask(ctx *gin.Context) {
	// Get task ID
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get task title
    title, exist := ctx.GetPostForm("title")
    if !exist {
        Error(http.StatusBadRequest, "No title is given")(ctx)
        return
    }
    // Get description
    description, _ := ctx.GetPostForm("description")
    // Get task priority
    priority, exist := ctx.GetPostForm("priority")
    if !exist {
        Error(http.StatusBadRequest, "No priority is given")(ctx)
        return
    }
    // Get task deadline
    deadline, exist := ctx.GetPostForm("deadline")
    if !exist {
        Error(http.StatusBadRequest, "No deadline is given")(ctx)
        return
    }
    // Get task tag
    tag_str, exist := ctx.GetPostForm("tag")
    tag := false
	if exist {
		tag, err = strconv.ParseBool(tag_str)
	    if err != nil {
            Error(http.StatusBadRequest, err.Error())(ctx)
            return
        }
	}
    // Get task category
    category, exist := ctx.GetPostForm("category")
    if !exist {
        Error(http.StatusBadRequest, "No category is given")(ctx)
        return
    }
	// Get task is_done
	is_done_str, exist := ctx.GetPostForm("is_done")
	if !exist {
		Error(http.StatusBadRequest, "No is_done is given")(ctx)
		return
	}
	is_done, err := strconv.ParseBool(is_done_str)
	if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get task share
    share_str, exist := ctx.GetPostForm("share")
    share := false
	if exist {
		share, err = strconv.ParseBool(share_str)
	    if err != nil {
            Error(http.StatusBadRequest, err.Error())(ctx)
            return
        }
	}
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	// edit data with given title and is_done on DB
	_, err = db.Exec("UPDATE tasks SET title = ?, is_done = ?, explanation = ?, priority = ?, deadline = ?, tag = ?, category = ? WHERE id = ?", title, is_done, description, priority, deadline, tag, category, id)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

    // Get groupID
    var groupID uint64
    if share {//共有するならgroupIDをuserIDから取得
        // Get UserID
        userID := sessions.Default(ctx).Get("user")
        err = db.Get(&groupID, "SELECT group_id FROM belong WHERE user_id = ?", userID)
        if err != nil {
            Error(http.StatusBadRequest, err.Error())(ctx)
            return
        }
    } else { //共有しないならgroupIDは0
        groupID = 0
    }
    // shareテーブルに登録する
    _, err = db.Exec("UPDATE share SET group_id = ? WHERE task_id = ?", groupID, id)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	// Render status
    path := fmt.Sprintf("/task/%d", id)
    ctx.Redirect(http.StatusFound, path)
}


func DeleteTask(ctx *gin.Context) {
    // ID の取得
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Delete the task from DB
    _, err = db.Exec("DELETE FROM tasks WHERE id=?", id)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Redirect to /list
    ctx.Redirect(http.StatusFound, "/list/1")
}


//グループ共有しているタスクは閲覧だけ可能にする
func SecureTaskGroup(ctx *gin.Context) {
    // Get task ID
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
        return
    }

    // Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
		return
	}

    // Get user_id from ownership
    var user_id uint64 //タスクを作成したユーザのID
    err = db.Get(&user_id, "SELECT user_id FROM ownership WHERE task_id = ?", id)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
        return
    }

    //ユーザがグループに所属しているか確認
    USERID := sessions.Default(ctx).Get("user") //閲覧しようとしているユーザのID
    var check int
    err = db.Get(&check, "SELECT COUNT(*) FROM belong WHERE user_id = ?", USERID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
        return
    }
    var group_id uint64 //ユーザが所属するグループID
    var task_group_id uint64 //閲覧するタスクが紐づいているグループiD
    if check > 0 {//もし所属しているなら
        //ユーザが所属しているグループを取得
        err = db.Get(&group_id, "SELECT group_id FROM belong WHERE user_id = ?", USERID)
        if err != nil {
            Error(http.StatusBadRequest, err.Error())(ctx)
            ctx.Abort()
            return
        }
        //閲覧しようとするタスクが紐づいているグループを取得
        err = db.Get(&task_group_id, "SELECT group_id FROM share WHERE task_id = ?", id)
        if err != nil {
            Error(http.StatusBadRequest, err.Error())(ctx)
            ctx.Abort()
            return
        }
    } else { //もし所属していないならgroup_idとtask_group_idは異なる値にする
        group_id = 1
        task_group_id = 2
    }

    // Check id
    if USERID != user_id && group_id != task_group_id {
        ctx.Redirect(http.StatusFound, "/login")
        ctx.Abort()
    } else {
        ctx.Next()
    }
}

func SecureTask(ctx *gin.Context) {
    // Get task ID
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
        return
    }

    // Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
		return
	}

    // Get user_id from ownership
    var user_id uint64 //タスクを作成したユーザのID
    err = db.Get(&user_id, "SELECT user_id FROM ownership WHERE task_id = ?", id)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        ctx.Abort()
        return
    }

    // Check id
    if sessions.Default(ctx).Get("user") != user_id {
        ctx.Redirect(http.StatusFound, "/login")
        ctx.Abort()
    } else {
        ctx.Next()
    }
}

