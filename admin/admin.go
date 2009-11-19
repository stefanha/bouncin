package admin

import (
	"log";
	"plugins";
)

type admin struct {
}

func (*admin) Enable() bool {
	log.Stderrf("admin enabled\n");
	return true;
}

func (*admin) Disable() {
	log.Stderrf("admin disabled\n")
}

func init() {
	plugins.Register("admin", &admin{})
}
