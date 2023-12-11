let userName = "";
let privateChatSocket;
let roomChatSocket;
let privateChat = Object.create(null);
let roomChat = Object.create(null);

let selectedPrivateChat = "";
let selectedRoomChat = "";
/*
	Message       string
	SendingUser   string
	ReceivingUser string
	TimeStamp     time.Time
*/

function connect(e){
    if (userName != "") return;
    let userInput = document.querySelector("#UsernameInput");
    let userInputValue = userInput.value;
    userName = userInputValue;

    privateChatSocket = new WebSocket("ws://localhost:8080/chat/");
    privateChatSocket.onopen = (e) => {
        let msg = {
            "Message": "",
            "SendingUser": userName,
            "ReceivingUser": "",
        }
        privateChatSocket.send(JSON.stringify(msg));
    }
    privateChatSocket.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        if (msg["userList"] != undefined){
            let userListUL = document.querySelector("#chatList");
            let userListInp =  msg["userList"].split("\n");
            userListUL.textContent = '';
            for (let i = 0; i < userListInp.length - 1; i++){
                let newUser = document.createElement("li");
                newUser.appendChild(document.createTextNode(userListInp[i]));
                userListUL.appendChild(newUser);
            }
        }
        else if (msg["roomList"] != undefined){
            let roomListUL = document.querySelector("#roomList");
            let roomListInp =  msg["roomList"].split("\n");
            roomListUL.textContent = '';
            for (let i = 0; i < roomListInp.length - 1; i++){
                let newRoom = document.createElement("li");
                newRoom.appendChild(document.createTextNode(roomListInp[i]));
                roomListUL.appendChild(newRoom);
            }
        }
        else{
            let sendingUser = msg["SendingUser"];
            if (sendingUser != "" && sendingUser != null) {
                let sendingUserMessage = msg["Message"];
                updatePrivateChatMessages(sendingUser, sendingUser, sendingUserMessage);
                refreshPrivateChat();
            }
        }

    }
}

/*START OF PRIVATE CHAT METHODS*/
function refreshPrivateChat(){
    if (selectedPrivateChat != "") {
        document.querySelector("#privateChatWindow").innerHTML = privateChat[selectedPrivateChat];
    }
}
function selectPrivateChat(){
    let inp = document.querySelector("#receiverInput");
    selectedPrivateChat = inp.value;
    refreshPrivateChat();
}

function updatePrivateChatMessages(chat, user, msg){
    if (privateChat[chat] == undefined) {
        privateChat[chat] = `${user}: ${msg} \n`;
    }
    else{
        privateChat[chat] += `${user}: ${msg} \n`;
    }
}

function sendPrivateMessage(e){
    let msgText = document.querySelector("#privateChatTextArea").value;
    let msg = {
        "Message": msgText,
        "SendingUser": "",
        "ReceivingUser": selectedPrivateChat
    }

    privateChatSocket.send(JSON.stringify(msg));

    if (msgText != "" && msgText[0] != "!"){
        updatePrivateChatMessages(selectedPrivateChat, userName, msgText);
        refreshPrivateChat();
    }
                
}


/*START OF ROOM CHAT METHODS*/
function connectToRoomChat(){
    if (userName == "") return;
    if (selectedRoomChat == "") return;

    if (roomChatSocket != undefined) roomChatSocket.close();
    roomChatSocket = new WebSocket("ws://localhost:8080/rooms/"+selectedRoomChat);

    roomChatSocket.onopen = (e) => {
        let msg = {
                "Message": "",
                "UserName": userName
        }
        roomChatSocket.send(JSON.stringify(msg));
    }

    roomChatSocket.onmessage = (e) => {
        let msg = JSON.parse(e.data);
        let userName = msg["UserName"];
        let localRoomChatMsg = msg["Message"]
        updateRoomChatMessages(selectedRoomChat, userName, localRoomChatMsg);
        refreshRoomChat();
    }
}

function refreshRoomChat(){
    if (selectedRoomChat != "") {
        document.querySelector("#roomChatWindow").innerHTML = roomChat[selectedRoomChat];
    }
}

function selectRoomChat(){
    let inp = document.querySelector("#roomNameInput");
    selectedRoomChat = inp.value;

    if (selectedRoomChat != ""){
        connectToRoomChat();
    }

    refreshRoomChat();
}

function updateRoomChatMessages(chat, user, msg){
    if (roomChat[chat] == undefined) {
        roomChat[chat] = `${user}: ${msg} \n`;
    }
    else{
        roomChat[chat] += `${user}: ${msg} \n`;
    }
}

function sendRoomChatMessage(e){
    let msgText = document.querySelector("#roomChatTextArea").value;
    let msg = {
        "Message": msgText,
        "UserName": ""
    }

    roomChatSocket.send(JSON.stringify(msg));

    if (msgText != "" && msgText[0] != "!"){
        updateRoomChatMessages(selectedRoomChat, userName, msgText);
        refreshRoomChat();
    }
                
}