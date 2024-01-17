TARGET_DIR := /tmp/pix-console
VERSION := 0.4
RELEASE := 4
STUNE_DIR := edward

# 產生執行檔
prepare:
	
build:
	go build
	sed -i 's/^Version: .*/Version: $(VERSION)/' pix-console.spec
	sed -i 's/^Release: .*/Release: $(RELEASE)/' pix-console.spec
	rm -rf $(TARGET_DIR)
	mkdir $(TARGET_DIR)
	mkdir $(TARGET_DIR)/logs
	cp -rp docs $(TARGET_DIR)
	cp -rp static $(TARGET_DIR)
	cp -rp templates $(TARGET_DIR)
	cp -rp config $(TARGET_DIR)
	cp -rp rbac $(TARGET_DIR)
	cp pix-console $(TARGET_DIR)
	rpmbuild -bb pix-console.spec

# 清理目標檔案
clean:
	rm -rf $(TARGET_DIR)
	rm -f pix-console
	rm -f pix-console.tar
	rm -f logs/*
	rm -rf pix-console*.rpm
	rm __debug_*

# 清理目標檔案
upload:
	./stune-tool upload /root/rpmbuild/RPMS/x86_64/pix-console-$(VERSION)-$(RELEASE).x86_64.rpm $(STUNE_DIR)
	echo pix-console-$(VERSION)-$(RELEASE).x86_64.rpm > version.txt
	./stune-tool upload version.txt $(STUNE_DIR)
download:
	./stune-tool download pix-console-$(VERSION)-$(RELEASE).x86_64.rpm $(STUNE_DIR)

install:
	rpm -ivh pix-console-$(VERSION)-$(RELEASE).x86_64.rpm

update:
	rpm -Uvh pix-console-$(VERSION)-$(RELEASE).x86_64.rpm

showVersion:
	rpm -qa|grep pix

journalStatus:
	journalctl -u pix-console -f

stop: 
	systemctl stop pix-console

start: 
	systemctl start pix-console
test: build upload download update
