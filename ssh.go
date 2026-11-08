package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SshUser struct {
	Username  string `json:"username"`
	UID       int    `json:"uid"`
	Shell     string `json:"shell"`
	Password  string `json:"password"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

func listSshUsers() []SshUser {
	users := []SshUser{}
	db := loadDB()
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return users
	}
	defer f.Close()

	skip := map[string]bool{"nobody": true, "ubuntu": true, "sshd": true, "dropbear": true, "stunnel": true}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), ":")
		if len(parts) < 7 {
			continue
		}
		uid, _ := strconv.Atoi(parts[2])
		if uid < 1000 || skip[parts[0]] {
			continue
		}
		u := SshUser{Username: parts[0], UID: uid, Shell: parts[6]}
		if meta, ok := db[parts[0]]; ok {
			u.Password = meta.Password
			u.IP = meta.IP
			u.UserAgent = meta.UserAgent
		} else {
			u.Password = "-"
			u.IP = "Unknown"
			u.UserAgent = "Unknown"
		}
		users = append(users, u)
	}
	return users
}

func addSshUser(username, password, ip, ua, activeHost string) (bool, string) {
	if username == "" || password == "" {
		return false, "Username dan password wajib diisi!"
	}
	if !validUsername(username) || !validPassword(password) {
		return false, "Username/Password mengandung karakter ilegal!"
	}
	if err := exec.Command("useradd", "-m", "-s", "/bin/bash", username).Run(); err != nil {
		return false, "Gagal membuat user. Username mungkin sudah terpakai."
	}
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", username, password))
	if err := cmd.Run(); err != nil {
		return false, "Gagal set password user"
	}
	db := loadDB()
	db[username] = UserMeta{Password: password, IP: ip, UserAgent: ua}
	saveDB(db)

	msg := fmt.Sprintf(
		"================================\n"+
			"       SSH ACCOUNT CREATED      \n"+
			"================================\n"+
			"🔹 Host SSH  : %s\n"+
			"🔹 Port TLS  : 443\n"+
			"🔹 Port NTLS : 80\n"+
			"🔹 Username  : %s\n"+
			"🔹 Password  : %s\n"+
			"================================\n"+
			" powered by : gogo\n"+
			"================================",
		activeHost, username, password)
	return true, msg
}

func deleteSshUser(username string) (bool, string) {
	if !validUsername(username) {
		return false, "Username ilegal!"
	}
	if err := exec.Command("userdel", "-r", username).Run(); err != nil {
		return false, "Gagal menghapus user."
	}
	db := loadDB()
	delete(db, username)
	saveDB(db)
	return true, "User " + username + " berhasil dihapus!"
}

func validUsername(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			return false
		}
	}
	return true
}

func validPassword(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '_' || r == '@' || r == '.' || r == '-') {
			return false
		}
	}
	return true
}
