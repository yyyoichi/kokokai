package user

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func testLoadEnv() {
	currentDir, _ := os.Getwd()
	envPath := strings.ReplaceAll(filepath.Join(currentDir, "../../.env"), "\\", "/")
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func TestCreateUser(t *testing.T) {
	testLoadEnv()
	u := &User{Id: "yyyoichi", Pass: "pa55w0rd"}
	err := u.Create()
	if err != nil {
		t.Errorf(err.Error())
	}
	if u.Pk == 0 {
		t.Error("Pk is 0.")
	}
	if u.Id == "" {
		t.Error("Id in ''")
	}
	t.Logf("Pk: %d, Email: %s, Id: %s, Pass: %s", u.Pk, u.Email, u.Id, u.Pass)
	t.Log(u.CreateAt)
}
func TestGetUser(t *testing.T) {
	testLoadEnv()
	u := &User{Id: "yyyoichi", Pass: "pa55w0rd"}
	err := u.GetByPass()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(u)
}
func TestDeleteUser(t *testing.T) {
	testLoadEnv()
	u := &User{Id: "yyyoichi", Pass: "pa55w0rd"}
	err := u.GetByPass()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(u)
	err = u.Delete()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = u.GetByPass()
	if err != nil {
		t.Logf(err.Error())
	} else {
		t.Errorf("exits")
	}
}

func UserEcosystem(u *User, t *testing.T) func() {
	err := u.Create()
	if err != nil {
		t.Errorf(err.Error())
	}
	if u.Pk == 0 {
		t.Error("Pk is 0.")
	}
	if u.Id == "" {
		t.Error("Id is empty")
	}
	return func() {
		u.Delete()
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestGetById(t *testing.T) {
	testLoadEnv()
	u := &User{Id: "yyyoichi", Pass: "pa55w0rd"}
	delete := UserEcosystem(u, t)
	defer delete()
	u.Name = ""
	u.GetById()
	if u.Name == "" {
		t.Errorf("expeced some string, got=''")
	}
}

func TestUpdate(t *testing.T) {
	testLoadEnv()
	u := &User{Id: "yyyoichi", Pass: "pa55w0rd"}
	delete := UserEcosystem(u, t)
	defer delete()

	u.Name = "updatename1"
	if err := u.Update(); err != nil {
		t.Error(err)
	}
	u.GetById()
	if u.Name != "updatename1" {
		t.Errorf("expeced updatename1. but got=%s", u.Name)
	}

	u.Name = "updatename2"
	u.Email = "updateemail2"
	if err := u.Update(); err != nil {
		t.Error(err)
	}
	u.GetById()
	if u.Name != "updatename2" {
		t.Errorf("expeced updatename2. but got=%s", u.Name)
	}
	if u.Email != "updateemail2" {
		t.Errorf("expeced updateemail2. but got=%s", u.Email)
	}

	u.Name = ""
	u.Email = ""
	if err := u.Update(); err != nil {
		t.Error(err)
	}
	u.GetById()
	if u.Name != "" {
		t.Errorf("expeced ''. but got=%s", u.Name)
	}
	if u.Email != "" {
		t.Errorf("expeced ''. but got=%s", u.Email)
	}
}
