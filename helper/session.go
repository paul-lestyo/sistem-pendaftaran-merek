package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var sessionStore *session.Store

func init() {
	sessionStore = session.New()
}

func GetSession(c *fiber.Ctx, key string) interface{} {
	sess, err := sessionStore.Get(c)
	PanicIfError(err)

	return sess.Get(key)
}

func SetSession(c *fiber.Ctx, key string, value string) {
	sess, err := sessionStore.Get(c)
	PanicIfError(err)

	sess.Set(key, value)
	err = sess.Save()
	PanicIfError(err)
}

func DeleteSession(c *fiber.Ctx, key string) {
	sess, err := sessionStore.Get(c)
	PanicIfError(err)
	sess.Delete(key)

	err = sess.Save()
	PanicIfError(err)
}

func DestroySession(c *fiber.Ctx) {
	sess, err := sessionStore.Get(c)
	PanicIfError(err)

	err = sess.Destroy()
	PanicIfError(err)
}
