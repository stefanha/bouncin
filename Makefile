include $(GOROOT)/src/Make.$(GOARCH)

TARG=bouncin
GOFILES=\
		irc.go\

include $(GOROOT)/src/Make.pkg
