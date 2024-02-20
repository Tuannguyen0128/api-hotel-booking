package service

import (
	"context"
	"sync"

	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/utils/emailutil"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence/dbmongo"
	"api-hotel-booking/internal/app/thirdparty/database"
	"api-hotel-booking/internal/app/thirdparty/mongo"
)

type Adapters struct {
	userRepo    persistence.UserDB
	sessionRepo persistence.SessionDB
	emailUtil   emailutil.EmailUtil
}

var (
	lock             = &sync.Mutex{}
	adaptersInit sync.Once
	adaptersInst *Adapters
	databaseInstance database.Manager
)

func NewAdapters() *Adapters {
	adaptersInit.Do(initMySQLAdapters)
	return adaptersInst
}

func initAdapters() {

	// Initiate mongo repo
	log.Info("Initializing connection to mongo: %+v", app.CFG.Mongo.DialInfo.Addrs)
	mgClient, err := mongo.NewMongoClient(app.CFG.Mongo)
	if err != nil {
		log.Panic("failed to create mongo client: %v<br>%v", err, app.CFG.Mongo)
	}
	var userRepo = dbmongo.NewUserRepo(app.CFG.Mongo.DialInfo.Database, mgClient)
	var sessionRepo = dbmongo.NewSessionRepo(app.CFG.Mongo.DialInfo.Database, mgClient)

	log.Info("Connection to mongo established successfully.")

	emailUtil := emailutil.NewEmailService(app.CFG.Email)

	adaptersInst = &Adapters{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		emailUtil:   emailUtil,
	}
}

func initMySQLAdapters() {
	if databaseInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		log.Info("Creating single instance of type database.Manager")


		conf := &database.Config{
			DSN:  app.CFG.MySQL.DSN(),
			Name: app.CFG.MySQL.Type,
		}

		mgr, pErr := database.NewMysqlService(context.Background(), conf)
		if pErr != nil {
			log.Error("Could not read database file",pErr)
		} else {
			databaseInstance = mgr
			log.Info("Connected")
		}
	}
}

func (a *Adapters) UserRepo() persistence.UserDB {
	return a.userRepo
}

func (a *Adapters) SessionRepo() persistence.SessionDB {
	return a.sessionRepo
}

func (a *Adapters) EmailUtil() emailutil.EmailUtil {
	return a.emailUtil
}
