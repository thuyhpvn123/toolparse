
var output = document.getElementById("log-content");
var socket = new WebSocket("ws://localhost:2000/ws");
var socketActive = false;
var $createResult = document.getElementById('create-result');
var addCall ="6cb833c99ffae43652fe96f366daeb45a45948cb";
var prikeyCall = "18366665a92ec88e3f1d765713fb6c23177f6efd4b226200cb84c1f216f2646f";
var kventureAdd = "0xf28B2FDA437D4717915FfDA6b858B6A29fBAE74f"
console.log("Imported");
// * Websocket
// Connect to server successfully

socket.onopen = (msg) => {
  socketActive = true;
  connectAWallet(addCall,prikeyCall)
// Your input data (a string)
};

// WS connection's closed
socket.onclose = (event) => {
  console.log("WS Connection is closed: ", event);
};

// WS connection having errors
socket.onerror = (error) => {
  console.log("Socket Error: ", error);
};
socket.onmessage = (msg) => {
  var data12 = JSON.parse(msg.data);
  output.innerHTML += "Server: " + msg.data + "\n";


}
var sendMessage = (msg) => {
  console.log(msg);
  socket.send(JSON.stringify(msg));

};
//connect a wallet
var $connectAWallet = document.getElementById('connectAWallet');

$connectAWallet.addEventListener('submit', async(e) => {
  
      e.preventDefault()
      console.log("connect-A-wallet")
      var flag=1,address,privateKey
      address = $('#address').val();
      privateKey = $('#privateKey').val()
      if( address ==''|| privateKey ==''){
        flag=0
        $('.error_getKey').html("Please type player address")
      }else{
        $('.error_getKey').html("")
      }
      if(flag==1  ){
        try{
          connectAWallet(address,privateKey)
        }catch{
          console.log(e)
          $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
        }
      }
  })
var connectAWallet =(address,privateKey)=>{
  var ms ={
    "address":address,
    "privateKey":privateKey
  }
  var setMsg = {
    command: "connect-wallet",
    value: ms,  
  }
  sendMessage(setMsg);
}
//call one to one
var $call = document.getElementById('call');

$call.addEventListener('click', async(e) => {
  
      e.preventDefault()
        try{
          call()
        }catch{
          console.log(e)
          $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
        }
  })
var call =()=>{
  var callms ={
        'from':   addCall,
        'to':   "0x79965F966Cdb14130127EEB9E8411869893ba26A",
        amount:            "",
        fee:             "1",
        'is-deploy':   false,
        'is-call':     true,
        'input': "0x062e9cd0",
        "function-name":"userViewProduct",
        'relatedAddresses':[],
        'contract':"product",
      }
  var messageForm = {
    command:"call",
    value: callms,
    };
  sendMessage(messageForm);
}

//call viewtreeMatric
var $call = document.getElementById('viewtreeMatrix');

$call.addEventListener('click', async(e) => {
  
      e.preventDefault()
        try{
          callViewtreeMatrix()
        }catch{
          console.log(e)
          $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
        }
  })
var callViewtreeMatrix =()=>{
  var callms ={
    'from':   addCall,
    'to':   kventureAdd,
    amount:            "",
    fee:             "1",
    'is-deploy':   false,
    'is-call':     true,
    'input': "0x6e83822b000000000000000000000000b96fae96b378145c1b111b2b9c8d7f8703e42487",
    "function-name":"viewtreeMatrix",
    'relatedAddresses':[],
    'contract':"kventure",
  }
  var messageForm = {
    command:"viewtreeMatrix",
    value: callms,
    };
  sendMessage(messageForm);
}
//call getCodeList
var $call = document.getElementById('getCodeList');

$call.addEventListener('click', async(e) => {
  
      e.preventDefault()
        try{
          getCodeList()
        }catch{
          console.log(e)
          $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
        }
})
function getCodeList(){
  fetch("../frontend/file/viewtreeMatrixCall.json")
  .then((response) => response.json())
  .then((json) =>{
    if(json.length>0){
      for (i=0;i<json.length;i++){
        console.log(i);

      }
    }
  });
}   
var callGetCodeList =(add)=>{
  inputGetCodeList = "0x134c0377000000000000000000000000" + add
  var callms ={
    'from':   addCall,
    'to':   "9eeB553fC6d5D6f9b35a2728121317DCc2eA6C38",
    amount:            "",
    fee:             "1",
    'is-deploy':   false,
    'is-call':     true,
    'input': "0x134c0377000000000000000000000000853c63e6eef57b5037364f1e9276799379b84014",
    "function-name":"getCodeList",
    'relatedAddresses':[],
    'contract':"codepool",
  }
    var messageForm = {
    command:"getCodeList",
    value: callms,
    };
  sendMessage(messageForm);
}
function stripHexPrefix(address) {
  // Check if the address starts with '0x'
  if (address.startsWith('0x')) {
      // If it does, remove the '0x' prefix
      return address.slice(2);
  }
  // If it doesn't start with '0x', return the original address
  return address;
}
//call GetUserInfo
var $callGetUserInfo = document.getElementById('GetUserInfo');

$callGetUserInfo.addEventListener('click', async(e) => {
  
      e.preventDefault()
        try{
          getUserList()
        }catch{
          console.log(e)
          $createResult.innerHTML = `Ooops... there was an error while trying to set players`;
        }
})
function getUserList(){
  fetch("../frontend/file/viewtreeMatrixCall.json")
  .then((response) => response.json())
  .then((json) =>{
    console.log("hello",json);
    if(json.data.length>0){
      for (i=0;i<json.data.length;i++){
        console.log(json.data[i]);
        callGetUserInfo(stripHexPrefix(json.data[i]))
      }
    }
  });
}   
var callGetUserInfo =(add)=>{
  inputGetCodeList = "0x2c9387f6000000000000000000000000" + add
  var callms ={
    'from':   addCall,
    'to':   kventureAdd,
    amount:            "",
    fee:             "1",
    'is-deploy':   false,
    'is-call':     true,
    'input': inputGetCodeList,
    "function-name":"GetUserInfo",
    'relatedAddresses':[],
    'contract':"kventure",
  }
    var messageForm = {
    command:"getUserInfo",
    value: callms,
    };
  sendMessage(messageForm);
}
