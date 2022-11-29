/* placeholder file for JavaScript */
const confirm_task_delete = (id) => {
    if(window.confirm(`Task ${id} を削除します．よろしいですか？`)) {
        location.href = `/task/delete/${id}`;
    }
}

const confirm_user_delete = () => {
    if(window.confirm(`アカウントを削除します．よろしいですか？`)) {
        location.href = `/user/delete`;
    }
}
 
const confirm_task_update = (id) => {
    // 練習問題 7-2
    if(window.confirm(`Task ${id} を更新します. よろしいですか？`)) {
        document.EditTask.submit();
    }
}

const confirm_user_update = () => {
    if(window.confirm(`アカウントを更新します. よろしいですか？`)) {
        document.EditUser.submit();
    }
}

const confirm_group_update = () => {
    if(window.confirm(`グループ情報を更新します. よろしいですか？`)) {
        document.EditGroup.submit();
    }
}

function countdown (time) {
    var now = new Date() - time;
    document.write(now);
}