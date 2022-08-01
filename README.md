踩雷记录：
1. 浏览器不能升级成ws协议， 使用curl进行websocket测试的时候需要以下前缀,如果多参数不能识别，添加转义符 "\"
curl --include \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Host: 127.0.0.1:3000" \
     --header "Origin: http://127.0.0.1:3000" \
     --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
     --header "Sec-WebSocket-Version: 13" \
     http://127.0.0.1:3000/ws?uid=3\&toUid=5

2. 可以使用http://www.websocket-test.com/ // postman工具 进行websocket通信测试
ws://127.0.0.1:3000/ws?uid=3\&toUid=5
{"type":1, "content": "the second message"}
3. https://www.ruanyifeng.com/blog/2017/05/websocket.html
有空了解一下websocket 不是用空 必须要看