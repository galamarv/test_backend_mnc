package controllers

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/galamarv/test_backend_mnc/models"
	"github.com/google/uuid"
)

type TransactionController struct {
	web.Controller
}

func (t *TransactionController) TopUp() {
	userId := t.Ctx.Input.GetData("user_id").(string)

	var input struct {
		Amount float64 `json:"amount"`
	}

	if err := t.ParseForm(&input); err != nil {
		t.CustomAbort(400, "Invalid input")
	}

	o := orm.NewOrm()
	user := models.User{Id: userId}
	if err := o.Read(&user); err != nil {
		t.CustomAbort(400, "User not found")
	}

	balanceBefore := user.Balance
	user.Balance += input.Amount

	tx, err := o.Begin()
	if err != nil {
		t.CustomAbort(500, "Error starting transaction")
	}

	transaction := models.Transaction{
		Id:            uuid.New().String(), // Convert UUID to string
		UserId:        user.Id,
		Type:          "CREDIT",
		Amount:        input.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}

	_, err = o.Update(&user)
	if err != nil {
		tx.Rollback()
		t.CustomAbort(500, "Error updating user balance")
	}

	_, err = o.Insert(&transaction)
	if err != nil {
		tx.Rollback()
		t.CustomAbort(500, "Error inserting transaction record")
	}

	err = tx.Commit()
	if err != nil {
		t.CustomAbort(500, "Error committing transaction")
	}

	t.Data["json"] = map[string]interface{}{
		"status": "SUCCESS",
		"result": transaction,
	}
	t.ServeJSON()
}

func (t *TransactionController) Pay() {
	userId := t.Ctx.Input.GetData("user_id").(string)

	var input struct {
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks"`
	}

	if err := t.ParseForm(&input); err != nil {
		t.CustomAbort(400, "Invalid input")
	}

	o := orm.NewOrm()
	user := models.User{Id: userId}
	if err := o.Read(&user); err != nil {
		t.CustomAbort(400, "User not found")
	}

	if user.Balance < input.Amount {
		t.CustomAbort(400, "Balance is not enough")
	}

	balanceBefore := user.Balance
	user.Balance -= input.Amount

	tx, err := o.Begin()
	if err != nil {
		t.CustomAbort(500, "Error starting transaction")
	}

	transaction := models.Transaction{
		Id:            uuid.New().String(), // Convert UUID to string
		UserId:        user.Id,
		Type:          "DEBIT",
		Amount:        input.Amount,
		Remarks:       input.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}

	_, err = o.Update(&user)
	if err != nil {
		tx.Rollback()
		t.CustomAbort(500, "Error updating user balance")
	}

	_, err = o.Insert(&transaction)
	if err != nil {
		tx.Rollback()
		t.CustomAbort(500, "Error inserting transaction record")
	}

	err = tx.Commit()
	if err != nil {
		t.CustomAbort(500, "Error committing transaction")
	}

	t.Data["json"] = map[string]interface{}{
		"status": "SUCCESS",
		"result": transaction,
	}
	t.ServeJSON()
}

func (t *TransactionController) Transfer() {
	userId := t.Ctx.Input.GetData("user_id").(string)

	var input struct {
		TargetUser string  `json:"target_user"`
		Amount     float64 `json:"amount"`
		Remarks    string  `json:"remarks"`
	}

	if err := t.ParseForm(&input); err != nil {
		t.CustomAbort(400, "Invalid input")
	}

	o := orm.NewOrm()
	user := models.User{Id: userId}
	if err := o.Read(&user); err != nil {
		t.CustomAbort(400, "User not found")
	}

	targetUser := models.User{Id: input.TargetUser}
	if err := o.Read(&targetUser); err != nil {
		t.CustomAbort(400, "Target user not found")
	}

	if user.Balance < input.Amount {
		t.CustomAbort(400, "Balance is not enough")
	}

	balanceBefore := user.Balance
	user.Balance -= input.Amount
	targetBalanceBefore := targetUser.Balance
	targetUser.Balance += input.Amount

	tx, err := o.Begin()
	if err != nil {
		t.CustomAbort(500, "Error starting transaction")
	}

	transferTransaction := models.Transaction{
		Id:            uuid.New().String(), // Convert UUID to string
		UserId:        user.Id,
		Type:          "DEBIT",
		Amount:        input.Amount,
		Remarks:       input.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}

	receiveTransaction := models.Transaction{
		Id:            uuid.New().String(), // Convert UUID to string
		UserId:        targetUser.Id,
		Type:          "CREDIT",
		Amount:        input.Amount,
		Remarks:       input.Remarks,
		BalanceBefore: targetBalanceBefore,
		BalanceAfter:  targetUser.Balance,
		CreatedDate:   time.Now(),
	}

	_, err = o.Update(&user)
	_, err = o.Update(&targetUser)
	if err != nil {
		tx.Rollback()
		t.CustomAbort(500, "Error updating balances")
	}

	_, err = o.Insert(&transferTransaction)
	_, err = o.Insert(&receiveTransaction)
	if err != nil {
		tx.Rollback()
		t.CustomAbort(500, "Error inserting transaction record")
	}

	err = tx.Commit()
	if err != nil {
		t.CustomAbort(500, "Error committing transaction")
	}

	t.Data["json"] = map[string]interface{}{
		"status": "SUCCESS",
		"result": transferTransaction,
	}
	t.ServeJSON()
}

func (t *TransactionController) TransactionReport() {
	userId := t.Ctx.Input.GetData("user_id").(string)

	o := orm.NewOrm()
	var transactions []models.Transaction
	_, err := o.QueryTable("transaction").Filter("UserId", userId).All(&transactions)
	if err != nil {
		t.CustomAbort(500, "Error retrieving transactions")
	}

	t.Data["json"] = map[string]interface{}{
		"status": "SUCCESS",
		"result": transactions,
	}
	t.ServeJSON()
}
