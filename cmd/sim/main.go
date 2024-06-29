package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"social-network-service/internal/api/account"
	"social-network-service/internal/api/common"
	"social-network-service/internal/api/post"
	"strconv"
	"strings"
	"sync"
	"time"
)

const addr = "http://localhost:8080"

var (
	maleNames   []string
	femaleNames []string
	surnames    []string
	cities      []string
	posts       []string
)

var users = []string{}

var friends = make(map[string][]string)

var mutex = sync.RWMutex{}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Println("[COUNT] [DIR]")
		fmt.Println("")
		fmt.Println("Where:")
		fmt.Println("")
		fmt.Println("- COUNT is number of users.")
		fmt.Println("- DIR is a directory that contains the following files: surnames.txt, male-names.txt, female-names.txt, cities.txt, posts.txt.")
		return
	}

	count, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println("The first argument should be a number")
		os.Exit(1)
	}

	err = prepare(os.Args[2])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go start(wg)
	}

	wg.Wait()
}

func prepare(path string) error {
	entries, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	f := map[string]struct{}{
		"cities.txt":       {},
		"female-names.txt": {},
		"male-names.txt":   {},
		"posts.txt":        {},
		"surnames.txt":     {},
	}

	m := map[string]string{
		"cities.txt":       "",
		"female-names.txt": "",
		"male-names.txt":   "",
		"posts.txt":        "",
		"surnames.txt":     "",
	}

	for _, entry := range entries {
		_, found := m[entry.Name()]

		if !found {
			continue
		}

		if entry.IsDir() {
			continue
		}

		absPath, err := filepath.Abs(filepath.Join(path, entry.Name()))

		if err != nil {
			return err
		}

		m[entry.Name()] = absPath
		delete(f, entry.Name())
	}

	if len(f) != 0 {
		keys := make([]string, 0, len(f))
		for k := range f {
			keys = append(keys, k)
		}

		return fmt.Errorf("failed to open the following files: %v", keys)
	}

	cities, err = getFromFile(m["cities.txt"])

	if err != nil {
		return err
	}

	maleNames, err = getFromFile(m["male-names.txt"])

	if err != nil {
		return err
	}

	femaleNames, err = getFromFile(m["female-names.txt"])

	if err != nil {
		return err
	}

	surnames, err = getFromFile(m["surnames.txt"])

	if err != nil {
		return err
	}

	posts, err = getFromFile(m["posts.txt"])

	if err != nil {
		return err
	}

	return nil
}

func start(wg *sync.WaitGroup) {
	defer wg.Done()

	userId, err := registerUser()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addUserToStorage(userId)

	token, err := login(userId)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		op := rand.Intn(4)

		switch op {
		case 0:
			createPost(userId, token)
		case 1:
			readFeed(userId, token)
		case 2:
			addFriend(userId, token)
		case 3:
			removeFriend(userId, token)
		}

		delay, err := time.ParseDuration(fmt.Sprintf("%vms", rand.Intn(4000)+1000))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		time.Sleep(delay)
	}
}

func getFromFile(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	strs := strings.Split(string(bytes), "\n")

	return strs, nil
}

func registerUser() (string, error) {
	var name string

	gender := rand.Intn(2)

	switch gender {
	case 0:
		name = maleNames[rand.Intn(len(maleNames))]
	case 1:
		name = femaleNames[rand.Intn(len(femaleNames))]
	}

	surname := surnames[rand.Intn(len(surnames))]
	city := cities[rand.Intn(len(cities))]

	registerReq := account.RegisterRequest{
		FirstName:  name,
		SecondName: surname,
		Gender:     common.GenderMale,
		Birthdate:  "1990-01-01",
		Biography:  "Test biography",
		City:       city,
		Password:   "123456",
	}

	registerReqBytes, err := json.Marshal(registerReq)

	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(registerReqBytes)

	url := fmt.Sprintf("%s/user/register", addr)
	req, err := http.NewRequest(http.MethodPost, url, reader)

	if err != nil {
		return "", err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var registerResp account.RegisterResponse

	err = json.Unmarshal(respBytes, &registerResp)

	if err != nil {
		return "", err
	}

	log.Printf("user %v (%s %s) registered\n", registerResp.UserId, registerReq.FirstName, registerReq.SecondName)

	return registerResp.UserId, nil
}

func login(userId string) (string, error) {
	loginReq := account.LoginRequest{
		UserId:   userId,
		Password: "123456",
	}

	loginReqBytes, err := json.Marshal(loginReq)

	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(loginReqBytes)

	url := fmt.Sprintf("%s/login", addr)
	req, err := http.NewRequest(http.MethodPost, url, reader)

	if err != nil {
		return "", err
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var loginResp account.LoginResponse

	err = json.Unmarshal(respBytes, &loginResp)

	if err != nil {
		return "", err
	}

	log.Printf("user %v logged in\n", loginReq.UserId)

	return loginResp.Token, nil
}

func readFeed(userId, token string) error {
	url := fmt.Sprintf("%s/post/feed", addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Authorization", "Bearer "+token)

	q := req.URL.Query()
	q.Add("offset", "0")
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var readFeedResp post.ReadFeedResponse
	json.Unmarshal(respBytes, &readFeedResp)

	log.Printf("User %v read %v posts\n", userId, len(readFeedResp.Posts))

	return nil
}

func createPost(userId, token string) error {
	textIdx := rand.Intn(len(posts))
	text := posts[textIdx]

	createPostReq := post.CreatePostRequest{
		Text: text,
	}

	createPostReqBytes, err := json.Marshal(createPostReq)

	if err != nil {
		return err
	}

	reader := bytes.NewReader(createPostReqBytes)

	url := fmt.Sprintf("%s/post/create", addr)
	req, err := http.NewRequest(http.MethodPost, url, reader)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var createPostResp post.CreatePostResponse
	json.Unmarshal(respBytes, &createPostResp)

	log.Printf("User %v created post %v\n", userId, createPostResp.PostId)

	return nil
}

func addFriend(userId, token string) error {
	existingUsers := getUsersFromStorage()

	if len(existingUsers) <= 1 {
		return nil
	}

	var chosenUserId string

	for {
		chosenUserId = existingUsers[rand.Intn(len(existingUsers))]

		if chosenUserId == userId {
			continue
		}

		friendIds := getFriendsFromStorage(userId)

		if slices.Contains(friendIds, chosenUserId) {
			continue
		}

		break
	}

	url := fmt.Sprintf("%s/friend/set/%s", addr, chosenUserId)
	req, err := http.NewRequest(http.MethodPut, url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := http.DefaultClient
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	addFriendToStorage(userId, chosenUserId)

	log.Printf("User %v added user %v as friend\n", userId, chosenUserId)

	return nil
}

func removeFriend(userId, token string) error {
	friendUserIds := getFriendsFromStorage(userId)

	if len(friendUserIds) == 0 {
		return nil
	}

	friendUserId := friendUserIds[rand.Intn(len(friendUserIds))]

	url := fmt.Sprintf("%s/friend/delete/%s", addr, friendUserId)
	req, err := http.NewRequest(http.MethodPut, url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := http.DefaultClient
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	removeFriendFromStorage(userId, friendUserId)

	log.Printf("User %v removed friend %v\n", userId, friendUserId)

	return nil
}

func getFriendsFromStorage(userId string) []string {
	mutex.RLock()
	defer mutex.RUnlock()

	friendList := friends[userId]

	output := make([]string, len(friendList))
	copy(output, friendList)

	return friendList
}

func addFriendToStorage(userId, friendUserId string) {
	mutex.Lock()
	defer mutex.Unlock()

	friends[userId] = append(friends[userId], friendUserId)
}

func removeFriendFromStorage(userId, friendUserId string) {
	mutex.Lock()
	defer mutex.Unlock()

	userFriends := friends[userId]

	for i, existingFriendUserId := range userFriends {
		if existingFriendUserId != friendUserId {
			continue
		}

		friends[userId] = append(friends[userId][:i], friends[userId][i+1:]...)
	}
}

func getUsersFromStorage() []string {
	mutex.RLock()
	defer mutex.RUnlock()

	output := make([]string, len(users))
	copy(output, users)

	return output
}

func addUserToStorage(userId string) {
	mutex.Lock()
	defer mutex.Unlock()

	users = append(users, userId)
}
