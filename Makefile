include $(GOROOT)/src/Make.$(GOARCH)

TARG=bouncin
GOFILES=\
		network.go\

pkgdir=pkg/$(GOOS)_$(GOARCH)
PKGS=\
		irc\

all: $(TARG)
clean: clean-pkgs
	rm -f *.[$(OS)] $(TARG)

.PHONY: clean-pkgs
clean-pkgs:
	@for dir in $(PKGS); do \
		$(MAKE) -C $$dir clean; \
	done
	rm -rf pkg

_go_.$O: $(GOFILES) $(PKGS)
	$(GC) -I $(pkgdir) -o $@ $(GOFILES)

.PHONY: $(PKGS)
$(PKGS):
	$(MAKE) -C $@ install

$(TARG): _go_.$O
	$(LD) -L $(pkgdir) -o $@ _go_.$O
