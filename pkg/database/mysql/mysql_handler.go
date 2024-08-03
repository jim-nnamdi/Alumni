package mysql

import (
	"context"
	"database/sql"
	"io"
	"log"
	"time"

	"github.com/jim-nnamdi/jinx/pkg/model"
)

var (
	_ Database  = &mysqlDatabase{}
	_ io.Closer = &mysqlDatabase{}
)

type mysqlDatabase struct {
	createUser           *sql.Stmt
	checkUser            *sql.Stmt
	getUserByEmail       *sql.Stmt
	getBySessionKey      *sql.Stmt
	getUserPortfolios    *sql.Stmt
	getUserTransactions  *sql.Stmt
	createNewTransaction *sql.Stmt
	addNewForumPost      *sql.Stmt
	getSingleForumPost   *sql.Stmt
	getAllForums         *sql.Stmt
}

func NewMySQLDatabase(db *sql.DB) (*mysqlDatabase, error) {
	var (
		createUser           = "INSERT INTO users(username, password, email, degree, grad_year,current_job, phone, session_key,profile_picture,linkedin_profile,twitter_profile) VALUES(?,?,?,?,?,?,?,?,?,?,?);"
		checkUser            = "SELECT * FROM users where email = ? AND password=?;"
		getUserByEmail       = "SELECT * FROM users where email = ?;"
		getBySessionKey      = "SELECT * FROM users where session_key=?;"
		getUserPortfolios    = "SELECT * FROM portfolio_order WHERE `user_email` = ?;"
		getUserTransactions  = "SELECT * FROM transactions WHERE `user_email` = ?;"
		createNewTransaction = "INSERT INTO transactions(from_user_id,from_user_email, to_user_id, to_user_email,type,created_at,updated_at,amount,user_email) VALUES(?,?,?,?,?,?,?,?);"
		addNewForumPost      = "INSERT INTO forums(title, description, author, slug, created_at, updated_at) VALUES (?,?,?,?,?,?)"
		getSingleForumPost   = "SELECT * FROM forums WHERE `slug` = ?;"
		getAllForums         = "SELECT * FROM forums"
		database             = &mysqlDatabase{}
		err                  error
	)
	if database.createUser, err = db.Prepare(createUser); err != nil {
		return nil, err
	}
	if database.checkUser, err = db.Prepare(checkUser); err != nil {
		return nil, err
	}
	if database.getUserByEmail, err = db.Prepare(getUserByEmail); err != nil {
		return nil, err
	}
	if database.getBySessionKey, err = db.Prepare(getBySessionKey); err != nil {
		return nil, err
	}
	if database.getUserPortfolios, err = db.Prepare(getUserPortfolios); err != nil {
		return nil, err
	}
	if database.getUserTransactions, err = db.Prepare(getUserTransactions); err != nil {
		return nil, err
	}
	if database.createNewTransaction, err = db.Prepare(createNewTransaction); err != nil {
		return nil, err
	}
	if database.addNewForumPost, err = db.Prepare(addNewForumPost); err != nil {
		return nil, err
	}
	if database.getSingleForumPost, err = db.Prepare(getSingleForumPost); err != nil {
		return nil, err
	}
	if database.getAllForums, err = db.Prepare(getAllForums); err != nil {
		return nil, err
	}
	return database, nil
}

func (db *mysqlDatabase) CreateUser(ctx context.Context, username string, password string, email string, degree string, gradyear string, currentjob string, phone string, sessionkey string, profilepicture string, linkedinprofile string, twitterprofile string) (bool, error) {
	userQuery, err := db.createUser.ExecContext(ctx, username, password, email, degree, gradyear, currentjob, phone, sessionkey, profilepicture, linkedinprofile, twitterprofile)
	if err != nil {
		return false, err
	}
	lastId, err := userQuery.LastInsertId()
	if err != nil {
		return false, err
	}
	if lastId == 0 || lastId < 1 {
		return false, err
	}
	return true, nil
}

func (db *mysqlDatabase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	getUserByEmail := db.getUserByEmail.QueryRowContext(ctx, email)
	err := getUserByEmail.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Degree, &user.GradYear, &user.CurrentJob, &user.Phone, &user.SessionKey, &user.ProfilePicture, &user.LinkedinProfile, &user.TwitterProfile)
	if err != nil {
		log.Println("get user by email", err)
		return nil, err
	}
	return user, nil
}

func (db *mysqlDatabase) CheckUser(ctx context.Context, email string, password string) (*model.User, error) {
	user := &model.User{}
	getUserByEmail := db.checkUser.QueryRowContext(ctx, email, password)
	err := getUserByEmail.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Degree, &user.GradYear, &user.CurrentJob, &user.Phone, &user.SessionKey, &user.ProfilePicture, &user.LinkedinProfile, &user.TwitterProfile)
	if err != nil {
		log.Println("checkuser", err)
		return nil, err
	}
	return user, nil
}

func (db *mysqlDatabase) GetBySessionKey(ctx context.Context, sessionkey string) (*model.User, error) {
	user := &model.User{}
	getBySessionKey := db.getBySessionKey.QueryRowContext(ctx, sessionkey)
	err := getBySessionKey.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Degree, &user.GradYear, &user.CurrentJob, &user.Phone, &user.SessionKey, &user.ProfilePicture, &user.LinkedinProfile, &user.TwitterProfile)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *mysqlDatabase) GetUserPortfolio(ctx context.Context, user_email string) (*[]model.PortfolioOrder, error) {
	var portfolioOrders = []model.PortfolioOrder{}
	var portfolioOrder = model.PortfolioOrder{}
	getPortfolio, err := db.getUserPortfolios.QueryContext(ctx, user_email)
	if err != nil {
		return nil, err
	}
	for getPortfolio.Next() {
		err := getPortfolio.Scan(&portfolioOrder.Id, &portfolioOrder.Type, &portfolioOrder.Security, &portfolioOrder.Unit, &portfolioOrder.Status, &portfolioOrder.Cancelled, &portfolioOrder.UserID, &portfolioOrder.UserEmail, &portfolioOrder.CreatedAt, &portfolioOrder.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	portfolioOrders = append(portfolioOrders, portfolioOrder)
	return &portfolioOrders, nil
}

func (db *mysqlDatabase) GetUserTransactions(ctx context.Context, user_email string) (*[]model.Transaction, error) {
	var transaction = model.Transaction{}
	var transactions = []model.Transaction{}
	getTransactions, err := db.getUserTransactions.QueryContext(ctx, user_email)
	if err != nil {
		return nil, err
	}
	for getTransactions.Next() {
		err := getTransactions.Scan(&transaction.Id, &transaction.FromUserID, &transaction.FromUserEmail, &transaction.ToUserID, &transaction.ToUserEmail, &transaction.TransactionType, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.Amount, &transaction.UserEmail)
		if err != nil {
			return nil, err
		}
	}
	transactions = append(transactions, transaction)
	return &transactions, nil
}

func (db *mysqlDatabase) CreateNewTransaction(ctx context.Context, from_user int, from_user_email string, to_user int, to_user_email string, transactiontype string, created_at time.Time, updated_at time.Time, amount int, user_email string) (bool, error) {
	createNewTx, err := db.createNewTransaction.ExecContext(ctx, from_user, from_user_email, to_user, to_user_email, transactiontype, created_at, updated_at, amount, user_email)
	if err != nil {
		return false, err
	}
	lastInsert, err := createNewTx.LastInsertId()
	if err != nil {
		return false, err
	}
	if lastInsert <= 0 {
		return false, err
	}
	return true, nil
}

func (db *mysqlDatabase) AddNewForumPost(ctx context.Context, title string, description string, author string, slug string, created_at time.Time, updated_at time.Time) (bool, error) {
	createNewForum, err := db.addNewForumPost.ExecContext(ctx, title, description, author, slug, created_at, updated_at)
	if err != nil {
		return false, err
	}
	lastInsert, err := createNewForum.LastInsertId()
	if err != nil {
		return false, err
	}
	if lastInsert <= 0 {
		return false, err
	}
	return true, nil
}

func (db *mysqlDatabase) GetSingleForumPost(ctx context.Context, slug string) (*model.Forum, error) {
	forum := &model.Forum{}
	getForumBySlug := db.getSingleForumPost.QueryRowContext(ctx, slug)
	err := getForumBySlug.Scan(&forum.Id, &forum.Title, &forum.Description, &forum.Author, &forum.Slug, &forum.CreatedAt, &forum.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (db *mysqlDatabase) GetAllForums(ctx context.Context) (*[]model.Forum, error) {
	var forum = model.Forum{}
	var forums = []model.Forum{}
	getForums, err := db.getAllForums.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for getForums.Next() {
		err := getForums.Scan(&forum.Id, &forum.Title, &forum.Description, &forum.Author, &forum.Slug, &forum.CreatedAt, &forum.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	forums = append(forums, forum)
	return &forums, nil
}

func (db *mysqlDatabase) Close() error {
	db.createUser.Close()
	db.checkUser.Close()
	db.getBySessionKey.Close()
	db.getUserByEmail.Close()
	db.getUserPortfolios.Close()
	db.getUserTransactions.Close()
	db.createNewTransaction.Close()
	db.addNewForumPost.Close()
	db.getSingleForumPost.Close()
	return nil
}
