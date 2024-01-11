document.getElementById("createFileButton").addEventListener("click", function () {
    // 发起GET请求以触发文件创建
    fetch("/api/v1/createfile")
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



document.getElementById("checkfile").addEventListener("click", function () {
    // 发起GET请求以触发文件创建
    fetch("/api/v1/checkfile")
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 将结果输出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出现错误：" + data.error;
            } else if (data.error) {
                // 将错误消息输出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 将错误消息输出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "发生错误：" + error;
        });
});


document.getElementById("checkmemory").addEventListener("click", function () {
    // 发起GET请求以获取内存信息
    fetch("/api/v1/checkmemory")
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 将错误消息输出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出现错误：" + data.error;
            } else {
                // 将结果输出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 将错误消息输出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "发生错误：" + error;
        });
});

document.getElementById("startPixCompose").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/start_pix_compose", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});

document.getElementById("stopPixCompose").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/stop_pix_compose", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});

document.getElementById("uploadToStune").addEventListener("click", function () {

    fetch("/api/v1/uploadToStune", {
        method: 'POST',
    })
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





document.getElementById("recordPort").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/getpacket?status=true", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});

document.getElementById("StoprecordPort").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/getpacket?status=false", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});


document.getElementById("StartListen").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/listen?status=true", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});
document.getElementById("StopListen").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/listen?status=false", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});


document.getElementById("StopListen").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/listen?status=false", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});


document.getElementById("UpdateServer").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/updateServer", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});

document.getElementById("UpdateDocker").addEventListener("click", function () {
    // 發起POST請求以獲取內存信息
    fetch("/api/v1/updateDocker", {
        method: 'POST', // 將請求方法改為POST
        headers: {
            'Content-Type': 'application/json', // 指定請求頭為JSON格式
            // 在這裡可以添加其他請求頭信息
        },
        // 如果需要發送數據，可以在這裡添加請求體
        // body: JSON.stringify({ key: 'value' })
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // 將錯誤消息輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = "出現錯誤：" + data.error;
            } else {
                // 將結果輸出到 outputResult 段落元素
                document.getElementById("outputResult").textContent = JSON.stringify(data, null, 2);
            }
        })
        .catch(error => {
            // 將錯誤消息輸出到 outputResult 段落元素
            document.getElementById("outputResult").textContent = "發生錯誤：" + error;
        });
});