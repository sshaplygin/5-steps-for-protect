package main

import (
	"context"
	"log"

	"backend/internal"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	app, err := internal.NewApp(ctx)
	if err != nil {
		log.Panicln(err)
	}
	defer app.Close()

	log := app.GetLogger()

	if err = app.InitTables(); err != nil {
		log.Panic("init tables", zap.Error(err))
	}

	if err = app.Run(ctx); err != nil {
		log.Panic("run app", zap.Error(err))
	}
}
