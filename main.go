package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/robfig/cron/v3"
	"go-worker/worker"
	"log"
	"time"
)

func main() {
	app := fiber.New()

	//middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(csrf.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "key"
		},
	}))
	app.Use(cors.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// cron
	cronjob := cron.New(cron.WithSeconds())

	_, _ = cronjob.AddFunc("*/10 * * * * *", func() {
		isRunning := worker.IsWorkerRunning()
		log.Printf("Is worker running: %v \n", isRunning)

		if !isRunning {
			emails1 := []string{
				"zonzon@gmail.com",
			}
			jQueue := worker.NewJobQueue(1)
			jQueue.Start()

			for _, email := range emails1 {
				s := worker.Sender{
					Email: email,
				}
				jQueue.JobRunning <- s
			}

			jQueue.Stop()
		}
	})

	cronjob.Start()

	log.Fatal(app.Listen(":3000"))
}
