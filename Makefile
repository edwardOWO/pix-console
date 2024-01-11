TARGET_DIR := /tmp/pix-console

# 產生執行檔
build:
	go build

# 清理目標檔案
clean:
	rm -rf $(TARGET_DIR)
	rm -f pix-console
	rm -f pix-console.tar
	rm -f logs/*


# 清理目標檔案
upload:
	rm -rf $(TARGET_DIR)
	mkdir $(TARGET_DIR)
	mkdir $(TARGET_DIR)/logs
	cp -rp docs $(TARGET_DIR)
	cp -rp static $(TARGET_DIR)
	cp -rp templates $(TARGET_DIR)
	cp -rp config $(TARGET_DIR)
	cp -rp rbac $(TARGET_DIR)
	cp pix-console $(TARGET_DIR)
	tar cvf pix-console.tar $(TARGET_DIR)
    ./stune-tool upload pix-console edward