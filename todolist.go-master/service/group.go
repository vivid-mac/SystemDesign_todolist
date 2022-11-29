package service
 
import (
    "net/http"
	"fmt"
	"strconv"
    "github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	database "todolist.go/db"
)

func NewGroupForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form_new_group.html", gin.H{"Title": "Register group"})
}

func RegisterGroup(ctx *gin.Context) {
	//フォームデータ受け取り
	groupname := ctx.PostForm("groupname")
	// DB 接続
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	// 重複チェック
    var duplicate int
    err = db.Get(&duplicate, "SELECT COUNT(*) FROM groups WHERE name=?", groupname)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    if duplicate > 0 {
        ctx.HTML(http.StatusBadRequest, "form_new_group.html", gin.H{"Title": "Register group", "Error": "Groupname is already taken", "Groupname": groupname})
        return
    }

	// DB への保存
    result, err := db.Exec("INSERT INTO groups(name) VALUES (?)", groupname)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	GroupID, err := result.LastInsertId()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	path := fmt.Sprintf("/group/new/complete/%d", GroupID)
	ctx.Redirect(http.StatusFound, path)
}

func Complete(ctx *gin.Context) {
	// Get Group ID
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }

	// DB 接続
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	// Get a group with given ID
	var group database.Group
	err = db.Get(&group, "SELECT * FROM groups WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	ctx.HTML(http.StatusOK, "group_complete.html", gin.H{"Title": "Register group complete", "Group": group})
}

func ShowGroup(ctx *gin.Context) {
	UserID := sessions.Default(ctx).Get("user")
    // Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

    // Get a group with given UserID
    var group database.Group
    err = db.Get(&group, "SELECT id, name FROM groups INNER JOIN belong ON group_id = id WHERE user_id = ?", UserID)
    if err != nil {
        ctx.HTML(http.StatusOK, "group.html", gin.H{"Group": nil})
		return
    }

	// Getusernames
	var usernames []string
	err = db.Select(&usernames, "SELECT name FROM users INNER JOIN belong ON user_id = id WHERE group_id = ?", group.ID)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}


    ctx.HTML(http.StatusOK, "group.html", gin.H{"Group": group, "Usernames": usernames})
}

func Belong(ctx *gin.Context) {
	groupname := ctx.PostForm("groupname")
	UserID := sessions.Default(ctx).Get("user")

	db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	// 最初からエラー文が出ないようにする
	if groupname == "" {
		ctx.HTML(http.StatusBadRequest, "belong.html", gin.H{"Title": "Belong", "Groupname": groupname, "Error": ""})
		return
	}

	// グループの取得
    var group database.Group
    err = db.Get(&group, "SELECT id, name FROM groups WHERE name = ?", groupname)
    if err != nil {
        ctx.HTML(http.StatusBadRequest, "belong.html", gin.H{"Title": "Belong", "Groupname": groupname, "Error": "No such group"})
        return
    }
	fmt.Printf("groupIDID%d\n", group.ID)

	// ユーザとグループ紐付け
	_, err = db.Exec("INSERT INTO belong(group_id, user_id) VALUES (?, ?)", group.ID, UserID)
	if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	ctx.Redirect(http.StatusFound, "/group")
}

func Leave(ctx *gin.Context) {
	userID := sessions.Default(ctx).Get("user")
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

	// Update group_id to 0 ユーザがグループを抜けたら所有するタスクのgroup_idを0にして無所属にする
	var task_ids []uint64 //ユーザが所有するタスクIDのスライス
	err = db.Select(&task_ids, "SELECT task_id FROM ownership WHERE user_id = ?", userID)
	if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	for _, id := range task_ids {
		_, err = db.Exec("UPDATE share SET group_id = 0 WHERE task_id = ?", id)
		if err != nil {
			Error(http.StatusBadRequest, err.Error())(ctx)
			return
		}
	}

	// Delete the belong from DB
	_, err = db.Exec("DELETE FROM belong WHERE user_id=?", userID)
    if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}

	// Redirect to "/"
    ctx.Redirect(http.StatusFound, "/")
}

func EditGroupForm(ctx *gin.Context) {
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
    var Group database.Group
    err = db.Get(&Group, "SELECT id, name FROM groups WHERE id=?", id)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Render edit form
    ctx.HTML(http.StatusOK, "form_edit_group.html",
        gin.H{"Title": fmt.Sprintf("Edit task %d", id), "Group": Group})
}

func UpdateGroup(ctx *gin.Context) {
	// Get task ID
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Get task title
    groupname, exist := ctx.GetPostForm("groupname")
    if !exist {
        Error(http.StatusBadRequest, "No name is given")(ctx)
        return
    }
    // Get DB connection
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
	// edit data with given title and is_done on DB
	_, err = db.Exec("UPDATE groups SET name = ? WHERE id = ?", groupname, id)
    if err != nil {
		Error(http.StatusBadRequest, err.Error())(ctx)
		return
	}
	// Render status
    ctx.Redirect(http.StatusFound, "/group")
}