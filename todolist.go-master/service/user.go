package service
 
import (
	"crypto/sha256"
	"encoding/hex"
    "net/http"
 
    "github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	database "todolist.go/db"
)
 
func NewUserForm(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "new_user_form.html", gin.H{"Title": "Register user"})
}

func hash(pw string) []byte {
    const salt = "todolist.go#"
    h := sha256.New()
    h.Write([]byte(salt))
    h.Write([]byte(pw))
    return h.Sum(nil)
}

//パスワードが半角英字,数字,記号の組み合わせかチェックする.
func password_check(pw string) bool{
    NumFlag := false
    CharFlag := false
    SymbolFlag := false
    for _, r := range pw {
        if '0' <= r && r <= '9' { //数字が含まれている場合
            NumFlag = true
        } else if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') { //半角英字が含まれている場合
            CharFlag = true
        } else if '!' <= r && r <= '~' { //記号が含まれている場合
            SymbolFlag = true
        } else { //それ以外が含まれている場合
            return false
        }
    }
    if NumFlag && CharFlag && SymbolFlag {
        return true
    } else {
        return false
    }
}

func RegisterUser(ctx *gin.Context) {
    // フォームデータの受け取り
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
	checkpass := ctx.PostForm("checkpass")
    switch {
    case username == "":
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Usernane is not provided", "Username": username, "Password": password, "Checkpass": checkpass})
		return
    case password == "":
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Password is not provided", "Username": username, "Password": password, "Checkpass": checkpass})
		return
	case checkpass == "":
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Checkpass is not provided", "Username": username, "Password": password, "Checkpass": checkpass})
		return
	case len(password) < 8:
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Use 8 or more characters", "Username": username, "Password": password, "Checkpass": checkpass})
		return
    case !password_check(password):
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Use password with only a mix of letters, numbers and symbols", "Username": username, "Password": password, "Checkpass": checkpass})
		return
	case password != checkpass:
		ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Those passwords didn't match", "Username": username, "Password": password, "Checkpass": checkpass})
		return
    }


    // DB 接続
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
 
    // 重複チェック
    var duplicate int
    err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    if duplicate > 0 {
        ctx.HTML(http.StatusBadRequest, "new_user_form.html", gin.H{"Title": "Register user", "Error": "Username is already taken", "Username": username, "Password": password, "Checkpass": checkpass})
        return
    }
 
    // DB への保存
    _, err = db.Exec("INSERT INTO users(name, password) VALUES (?, ?)", username, hash(password))
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    
	ctx.Redirect(http.StatusFound, "/login")
}

const userkey = "user"
 
func Login(ctx *gin.Context) {
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
 
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	// 最初からエラー文が出ないようにする
	if username == "" && password == "" {
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": ""})
		return
	}
 
    // ユーザの取得
    var user database.User
    err = db.Get(&user, "SELECT id, name, password FROM users WHERE name = ?", username)
    if err != nil {
        ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": "No such user"})
        return
    }
 
    // パスワードの照合
    if hex.EncodeToString(user.Password) != hex.EncodeToString(hash(password)) {
        ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Title": "Login", "Username": username, "Error": "Incorrect password"})
        return
    }
 
    // セッションの保存
    session := sessions.Default(ctx)
    session.Set(userkey, user.ID)
    session.Save()
 
    ctx.Redirect(http.StatusFound, "/list/1")
}

func LoginCheck(ctx *gin.Context) {
    if sessions.Default(ctx).Get(userkey) == nil {
        ctx.Redirect(http.StatusFound, "/login")
        ctx.Abort()
    } else {
        ctx.Next()
    }
}

func Logout(ctx *gin.Context) {
    session := sessions.Default(ctx)
    session.Clear()
    session.Options(sessions.Options{MaxAge: -1})
    session.Save()
    ctx.Redirect(http.StatusFound, "/")
}

func ShowUser(ctx *gin.Context) {
    UserID := sessions.Default(ctx).Get(userkey)
    // Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

    // Get a username with given UserID
    var name string
    err = db.Get(&name, "SELECT name FROM users WHERE id=?", UserID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
		return
    }

    ctx.HTML(http.StatusOK, "user.html", gin.H{"ID": UserID, "UserName": name})
}

func EditUser(ctx *gin.Context) {
    UserID := sessions.Default(ctx).Get(userkey)
    // Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		Error(http.StatusInternalServerError, err.Error())(ctx)
		return
	}

    // Get a username with given UserID
    var name string
    err = db.Get(&name, "SELECT name FROM users WHERE id=?", UserID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
		return
    }

    ctx.HTML(http.StatusOK, "form_edit_user.html", gin.H{"ID": UserID, "Username": name})
}

func UpdateUser(ctx *gin.Context) {
    UserID := sessions.Default(ctx).Get(userkey)
    // フォームデータの受け取り
    username := ctx.PostForm("username")
    oldpassword := ctx.PostForm("oldpassword")
    newpassword := ctx.PostForm("newpassword")
	checkpass := ctx.PostForm("checkpass")
    switch {
    case username == "":
        ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Usernane is not provided", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
    case oldpassword == "":
        ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Old Password is not provided", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
    case newpassword == "":
        ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "New Password is not provided", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
	case checkpass == "":
		ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Checkpass is not provided", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
	case len(newpassword) < 8:
		ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Use 8 or more characters", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
    case !password_check(newpassword):
        ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Use password with only a mix of letters, numbers and symbols", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
	case newpassword != checkpass:
		ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Those passwords didn't match", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
		return
    }

    // DB 接続
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

    // ユーザの取得
    var user database.User
    err = db.Get(&user, "SELECT id, name, password FROM users WHERE id = ?", UserID)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
 
    // パスワードの照合
    if hex.EncodeToString(user.Password) != hex.EncodeToString(hash(oldpassword)) {
        ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Incorrect password", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
        return
    }
 
    // ユーザ名が新しくなっていた場合は重複チェック
    if user.Name != username {
        var duplicate int
        err = db.Get(&duplicate, "SELECT COUNT(*) FROM users WHERE name=?", username)
        if err != nil {
            Error(http.StatusInternalServerError, err.Error())(ctx)
            return
        }
        if duplicate > 0 {
            ctx.HTML(http.StatusBadRequest, "form_edit_user.html", gin.H{"Title": "Register user", "Error": "Username is already taken", "Username": username, "OldPassword": oldpassword, "NewPassword": newpassword, "Checkpass": checkpass})
            return
        }
    }

    // edit user with given username and password on DB
    _, err = db.Exec("UPDATE users SET name = ?, password = ? WHERE id = ?", username, hash(newpassword), UserID)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

	ctx.Redirect(http.StatusFound, "/user")
}

func DeleteUser(ctx *gin.Context) {
    UserID := sessions.Default(ctx).Get(userkey)
    db, err := database.GetConnection()
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Get target ownership
    var TaskID []uint64
    err = db.Select(&TaskID, "SELECT task_id FROM ownership WHERE user_id=?", UserID)
    if err != nil {
        Error(http.StatusBadRequest, err.Error())(ctx)
        return
    }
    // Delete target task    
    for id := range TaskID {
        _, err = db.Exec("DELETE FROM tasks WHERE id=?", id)
        if err != nil {
            Error(http.StatusInternalServerError, err.Error())(ctx)
            return
        }
    }
    // Delete target ownership
    _, err = db.Exec("DELETE FROM ownership WHERE user_id=?", UserID)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }
    // Delete target user
    _, err = db.Exec("DELETE FROM users WHERE id=?", UserID)
    if err != nil {
        Error(http.StatusInternalServerError, err.Error())(ctx)
        return
    }

    // Logout and Redirect to /
    ctx.Redirect(http.StatusFound, "/logout")
}