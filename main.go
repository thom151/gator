package main

import (
	"database/sql"
	"github.com/thom151/gator/internal/config"
	"github.com/thom151/gator/internal/database"
	"log"
	"os"
)

import _ "github.com/lib/pq"

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
		return
	}

	//opening connection to a database
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting database: %v ", err)
	}

	dbQueries := database.New(db)
	currState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		commandMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	cmdArgs := args[2:]

	cmd := command{
		name: args[1],
		arg:  cmdArgs,
	}

	err = cmds.run(currState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
