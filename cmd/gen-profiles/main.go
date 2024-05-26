package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"social-network-service/internal/admin"
	"social-network-service/internal/database"
	"social-network-service/internal/model"
	"social-network-service/internal/repository"
	"strings"
	"sync"
	"time"
)

var (
	accountCount = 100_000
	batchSize    = 10_000
)

const (
	maleNamesFilePath   = "../../scripts/male-names.txt"
	femaleNamesFilePath = "../../scripts/female-names.txt"
	surnamesFilePath    = "../../scripts/surnames.txt"
	citiesFilePath      = "../../scripts/cities.txt"
)

type batch struct {
	startIdx int
	length   int
}

var maleNames = getFromFile(maleNamesFilePath)
var femaleNames = getFromFile(femaleNamesFilePath)
var surnames = getFromFile(surnamesFilePath)
var cities = getFromFile(citiesFilePath)

func init() {
	flag.IntVar(&accountCount, "count", accountCount, "Number of accounts to be created")
	flag.IntVar(&batchSize, "batch-size", batchSize, "Size of batches")
	flag.Parse()
}

func main() {
	fmt.Println("Applied settings:")
	fmt.Println("count:", accountCount)
	fmt.Println("batch-size:", batchSize)
	fmt.Println()

	cfCfg := database.ConnectionFactoryConfig{
		MasterConnectionString: "host=localhost port=15432 user=postgres password=123 dbname=social_network_db",
		SyncConnectionString:   "host=localhost port=25432 user=postgres password=123 dbname=social_network_db",
		AsyncConnectionString:  "host=localhost port=35432 user=postgres password=123 dbname=social_network_db",
	}

	cf := database.NewConnectionFactory(cfCfg)

	tm := database.NewTransactionManager(cf)

	userRepoConfig := repository.UserRepositoryConfiguartion{
		UseAsyncReplicaForReadOperations: false,
	}

	userRepository := repository.NewUserRepository(userRepoConfig, cf)
	userAccountRepository := repository.NewUserAccountRepository(cf)

	appServiceConfig := &admin.AdminServiceConfiguration{
		UserRepository:        userRepository,
		UserAccountRepository: userAccountRepository,
		TransactionManager:    tm,
	}

	appService := admin.NewAdminService(appServiceConfig)

	generateUsers(appService)
}

func generateUsers(service *admin.AdminService) {
	curIdx := 0
	batches := []batch{}

	for curIdx < accountCount {
		if curIdx+batchSize < accountCount {
			batches = append(batches, batch{startIdx: curIdx, length: batchSize})
		} else {
			batches = append(batches, batch{startIdx: curIdx, length: accountCount - curIdx})
		}

		curIdx += batchSize
	}

	wg := sync.WaitGroup{}

	wg.Add(len(batches))

	for curBatchIdx, curBatch := range batches {
		go func(batch *batch, batchIdx int) {
			defer wg.Done()

			log.Printf("Batch %v: Creating users", batchIdx)

			cmds := make([]*model.RegisterUserCommand, batch.length)

			cmdIdx := 0

			for i := batch.startIdx; i < batch.startIdx+batch.length; i++ {
				isMale := generateGender()
				name := generateName(isMale)
				surname := generateSurname()
				birthdate := generateBirthdate()
				city := generateCity()

				var gender model.Gender

				if isMale {
					gender = model.GenderMale
				} else {
					gender = model.GenderFemale
				}

				cmd := &model.RegisterUserCommand{
					FirstName:  name,
					SecondName: surname,
					Gender:     gender,
					Birthdate:  birthdate,
					City:       city,
					Biography:  "Test",
					Password:   "123456",
				}

				cmds[cmdIdx] = cmd

				cmdIdx++
			}

			fmt.Printf("Batch %v: Registering %v users\n", batchIdx, len(cmds))

			err := service.RegisterUsers(context.Background(), cmds)

			if err != nil {
				panic(err)
			}

			fmt.Printf("Batch %v: Successfully registered %v users\n", batchIdx, len(cmds))
		}(&curBatch, curBatchIdx)
	}

	wg.Wait()

	fmt.Println("Mass user generation has been completed")
}

func getFromFile(path string) []string {
	_, execPath, _, ok := runtime.Caller(0)

	if !ok {
		panic("failed to get the path of executed file")
	}

	execDir := filepath.Dir(execPath)

	filePath := filepath.Join(execDir, path)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)

	if err != nil {
		panic("failed to open file")
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		panic("failed to read file")
	}

	strs := strings.Split(string(bytes), "\r\n")

	return strs
}

func generateGender() bool {
	x := rand.Intn(2)

	switch x {
	case 0:
		return true
	case 1:
		return false
	}

	panic("failed to generate gender due to wrong number")
}

func generateName(isMale bool) string {
	if isMale {
		nameCount := len(maleNames)
		randIdx := rand.Intn(nameCount)
		randomName := maleNames[randIdx]
		return randomName
	} else {
		nameCount := len(femaleNames)
		randIdx := rand.Intn(nameCount)
		randomName := femaleNames[randIdx]
		return randomName
	}
}

func generateSurname() string {
	surnameCount := len(surnames)
	randIdx := rand.Intn(surnameCount)
	randomSurname := surnames[randIdx]

	return randomSurname
}

func generateBirthdate() time.Time {
	yearOffset := rand.Intn(30)
	year := 2015 - yearOffset

	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1

	birthDate, err := time.Parse("2006-01-02", fmt.Sprintf("%v-%02d-%02d", year, month, day))

	if err != nil {
		panic("failed to parse birthdate")
	}

	return birthDate
}

func generateCity() string {
	cityCount := len(cities)
	randIdx := rand.Intn(cityCount)
	randomCity := cities[randIdx]

	return randomCity
}
