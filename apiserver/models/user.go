package models

import (
	"gogo-cmdb/apiserver/forms"
	"gogo-cmdb/apiserver/utils"
	"time"
)

type User struct {
	ID         int              `db:"id" json:"userId,string"`
	Name       string           `db:"username" json:"userName"`
	Password   string           `db:"password" json:"-"`
	UserType   utils.UserType   `db:"usertype" json:"userType,string"`
	UserStatus utils.UserStatus `db:"userstatus" json:"userStatus,string"`
}

type UserDetail struct {
	ID       int              `db:"-" json:"-"`
	UserId   int              `db:"-" json:"-"`
	Gender   int              `db:"gender" json:"gender,string"`
	Birthday *time.Time       `db:"birthday" json:"birthDay"`
	Tel      string           `db:"tel" json:"tel"`
	Email    string           `db:"email" json:"email"`
	Addr     string           `db:"addr" json:"addr"`
	Remark   string           `db:"remark" json:"remark"`
	Avatar   utils.NullString `db:"avatar" json:"avatar"`
}

type UserTimeStamp struct {
	ID          int        `db:"-" json:"-"`
	UserId      int        `db:"-" json:"-"`
	CreatedTime *time.Time `db:"createtime" json:"createdTime"`
	UpdatedTime *time.Time `db:"updatetime" json:"updatedTime"`
}

type UserFullInfo struct {
	User
	UserDetail
	UserTimeStamp
}

type UserManager struct{}

var DefalutUserManager *UserManager

func init() {
	DefalutUserManager = new(UserManager)
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (u *UserManager) GetUserList(page, size int) (int, []*User, *utils.Pagination, error) {
	sql_count := `
		select count(username) from users where userstatus !=?;
	`
	sql_users := `
		select id, username, password, usertype, userstatus from users where userStatus !=? order by id desc limit ?,?;
	`
	var total int
	err := db.Get(&total, sql_count, utils.Deleted)
	if err != nil {
		utils.Logger.Error(err)
		return 0, nil, nil, err
	}

	pagination, err := utils.NewPagination(total, page, size, 5)
	if err != nil {
		utils.Logger.Error(err)
		return 0, nil, nil, err
	}

	var userList []*User

	err = db.Select(&userList, sql_users, utils.Deleted, pagination.CurrentPage.Offset, pagination.CurrentPage.Limit)
	if err != nil {
		utils.Logger.Error(err)
		return 0, nil, nil, err
	}

	return total, userList, pagination, nil
}

func (u *UserManager) GetUserByName(name string) (*User, error) {
	sql_user := `
		select id, username, password, usertype, userstatus from users where username = ?;
	`
	//fmt.Println("name:", name)
	var user User
	err := db.Get(&user, sql_user, name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserManager) GetUserById(id int) (*User, error) {
	sql_user := `
		select id, username, password, usertype, userstatus from users where id = ?;
	`
	var user User
	err := db.Get(&user, sql_user, id)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return &user, nil
}

func (u *UserManager) GetUserDetailById(uid int) (*UserFullInfo, error) {
	sql_user := `
		select users.id, username, password, usertype, userstatus, 
		gender, birthday, tel, email, addr, remark, avatar, 
		createtime, updatetime 
		from users inner join user_details inner join user_timestamp on users.id = user_details.uid and users.id = user_timestamp.uid 
		where users.id = ?;
	`
	var user_full_info UserFullInfo

	err := db.Get(&user_full_info, sql_user, uid)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	return &user_full_info, nil
}

func (u *UserManager) GetUserDetailByName(uname string) (*UserFullInfo, error) {
	sql_user := `
		select users.id, username, password, usertype, userstatus, 
		gender, birthday, tel, email, addr, remark, avatar, 
		createtime, updatetime 
		from users inner join user_details inner join user_timestamp on users.id = user_details.uid and users.id = user_timestamp.uid 
		where users.username = ?;
	`
	var user_full_info UserFullInfo

	err := db.Get(&user_full_info, sql_user, uname)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	//fmt.Printf("%+v\n", details)
	return &user_full_info, nil
}

func (u *UserManager) CreateUser(urf *forms.UserRegisterForm) (int64, error) {

	createTime := time.Now()
	user_sql := `insert into users(usertype, userstatus, username, password) values(?,?,?,?);`
	detail_sql := `insert into user_details(uid, gender, tel, email, birthday, addr, remark) values(?,?,?,?,?,?,?);`
	timestamp_sql := `insert into user_timestamp(uid, createtime) values(?,?);`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return -1, err
	}
	defer tx.Rollback()
	result, err := tx.Exec(user_sql, utils.IntToUserType(urf.UserType), utils.IntToUserStatus(urf.UserStatus), urf.UserNmae, utils.Md5SaltPass(urf.Password, ""))
	if err != nil {
		utils.Logger.Error(err)
		return -1, err
	}

	uid, err := result.LastInsertId()
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return -1, err
	}
	_, err = tx.Exec(detail_sql, uid, urf.Gender, urf.Tel, urf.Email, urf.BirthDay, urf.Addr, urf.Remark)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return -1, err
	}
	_, err = tx.Exec(timestamp_sql, uid, createTime)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return -1, err
	}
	return uid, nil
}

func (u *UserManager) UpdateDetailById(uid int, uuf *forms.DetailUpdateForm) error {

	updateTime := time.Now()
	sql := `update user_details inner join user_timestamp on user_timestamp.uid = user_details.uid 
		set gender=?, tel=?, email=?, birthday=?, addr=?, remark=?, user_timestamp.updatetime=? where user_details.uid=?;`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sql, uuf.Gender, uuf.Tel, uuf.Email, uuf.BirthDay, uuf.Addr, uuf.Remark, &updateTime, uid)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (u *UserManager) UpdateAvatarById(url string, uid int) error {

	updateTime := time.Now()
	sql := `update user_details inner join user_timestamp on user_timestamp.uid = user_details.uid 
		set avatar=?, user_timestamp.updatetime=? where user_details.uid=?;`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sql, url, &updateTime, uid)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (u *UserManager) UpdateUserTypeById(uid int, utuf *forms.UserTypeUpdateForm) error {

	updateTime := time.Now()
	sql := `update users inner join user_timestamp on user_timestamp.uid = users.id 
		set usertype=?, user_timestamp.updatetime=? where users.id=?;`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sql, utils.IntToUserType(utuf.UserType), updateTime, uid)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (u *UserManager) UpdateUserStatusById(uid int, usuf *forms.UserStatusUpdateForm) error {

	updateTime := time.Now()
	sql := `update users inner join user_timestamp on user_timestamp.uid = users.id 
		set userstatus=?, user_timestamp.updatetime=? where users.id=?;`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sql, utils.IntToUserStatus(usuf.UserStatus), &updateTime, uid)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (u *UserManager) UpdatePasswordById(uid int, upuf *forms.PasswordUpdateForm) error {
	updateTime := time.Now()
	sql := `update users inner join user_timestamp on user_timestamp.uid = users.id 
		set password=?, user_timestamp.updatetime=? where users.id=?;`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sql, utils.Md5SaltPass(upuf.Password, ""), &updateTime, uid)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}

func (u *UserManager) UpdatePasswordByName(uName, pass string) error {
	updateTime := time.Now()
	sql := `update users inner join user_timestamp on user_timestamp.uid = users.id 
		set password=?, user_timestamp.updatetime=? where users.username=?;`

	tx, err := db.Begin()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(sql, utils.Md5SaltPass(pass, ""), &updateTime, uName)
	if err != nil {
		utils.Logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	return nil
}