document.getElementById("createFileButton").addEventListener("click", function() {
    // 发起GET请求以触发文件创建
    fetch("/createfile")
        .then(response => response.json())
        .then(data => {
            if (data.message) {
                alert(data.message);
            } else if (data.error) {
                alert("出现错误：" + data.error);
            }
        })
        .catch(error => {
            console.error("发生错误：", error);
        });
});

