TARGET_DIR := /tmp/pix-console
VERSION := 20240411
RELEASE := 15
STUNE_DIR := edward

# 產生執行檔
prepare:
	wget https://golang.org/dl/go1.21.2.linux-amd64.tar.gz
	tar -C /usr/local -xzf go1.21.2.linux-amd64.tar.gz
	rm go1.21.2*
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> ~/.bashrc
	source ~/.bashrc
	cd release-tool && go build -o stune-tool
	cp release-tool/stune-tool ./
	dnf --enablerepo=powertools install libpcap-devel -y
	sudo yum install rpm -y
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build
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
	rpmbuild --target=x86_64 -bb pix-console.spec
# create package
	./stune-tool download container.json edward
	mkdir -p $(VERSION)-$(RELEASE)
	cp container.json $(VERSION)-$(RELEASE)
	mv ~/rpmbuild/RPMS/x86_64/pix-console-$(VERSION)-$(RELEASE).x86_64.rpm $(VERSION)-$(RELEASE)
	tar cvf $(VERSION)-$(RELEASE).tar $(VERSION)-$(RELEASE)

# 清理目標檔案
clean:
	rm stune-tool 
	rm -rf $(TARGET_DIR)
	rm -f pix-console
	rm -f pix-console.tar
	rm -f logs/*
	rm -rf pix-console*.rpm
	rm -f stune-tool*
	rm -f __debug_*

# 清理目標檔案
upload:
	./stune-tool upload ~/rpmbuild/RPMS/x86_64/pix-console-$(VERSION)-$(RELEASE).x86_64.rpm $(STUNE_DIR)
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
