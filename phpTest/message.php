<?php
	function curlFun($uri,$data){
		$postdata = http_build_query($data);
		 $opts = array(
		 	'http' =>array(
                          'method'  => 'POST',
                          'header'  => 'Content-type: application/x-www-form-urlencoded',
                          'content' => $postdata
                     )
        );
        $context = stream_context_create($opts);
        $result = file_get_contents($uri, false, $context);
        return $result;
	}

	$ac = @$_GET["ac"];
	//消息发送
	if($ac == "postMsg"){
		$time = time();
		$param = array(
				"SendType"=>$_POST['SendType'],
				"Msg"=>$_POST['Msg'],
				"Topic"=>$_POST['Topic'],
				"Data"=> json_decode($_POST['data'],false),
			);
		
		$data = array(
			"sign"=>md5($time."ajsdlkfjdslk@sd~!"),
			"time"=>strval($time),
			"data"=>json_encode($param)
		);
		echo curlFun($_POST["sendToUrl"],$data);

	}else if ($ac=="subscribe"){
		$data = array(
			"sign"=>md5($_POST['topics']."ajsdlkfjdslk@sd~!".$_POST['uid']),
			"topics"=>$_POST['topics'],
			"uid"=>$_POST['uid']
		);
		echo curlFun($_POST["sendToUrl"],$data);

	}else{
?><!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
	<title>cmop Notify 测试工具</title>
	<link rel="stylesheet" href="//cdn.bootcss.com/bootstrap/3.3.5/css/bootstrap.min.css">
 </head>
 <body>
<div class="row">

<div class="col-md-8">
<div class="col-md-6">
<div class="panel panel-default">
  <div class="panel-heading">消息服务器基础配置</div>
  <div class="panel-body">
     <div class="form-group">
    <label for="exampleInputEmail1">消息服务器websocket地址</label>
    <input type="text" class="form-control" id="websocketUrl" placeholder="ws://" value="ws://c.com:9091/ws">
  </div>
  <div class="form-group">
    <label for="exampleInputPassword1">消息服务器http接口地址</label>
    <input type="text" class="form-control" id="apiUrl" placeholder="http://" value="http://c.com:9091">
  </div>


  </div>
</div>
</div>


<div class="col-md-6">
<div class="panel panel-default">
  <div class="panel-heading">用户配置</div>
  <div class="panel-body">
     <div class="form-group">
    <label for="exampleInputEmail1">用户标识</label>
    <input type="text" class="form-control" id="userid" placeholder="user1" value="testuser">
  </div>
  <div class="form-group">
    <label for="exampleInputPassword1">订阅话题</label>
    <input type="text" class="form-control" id="topics" placeholder="connectDesktop165241,connectDesktop165241" value="SNBNHTTX,TL0BCI78,UEMAVXNN,MDA2KG24,R5SDTXHF,71EUEBLL,KEKJ6G9K
    ">
  </div>

  </div>
</div>
</div>

<div class="col-md-12">
	
<div class="panel panel-default">
  <div class="panel-heading">连接</div>
  <div class="panel-body">
  	<div class="form-group">
       <button class="btn btn-primary" type="botton" onclick="connectWebSocket()" >连接消息服务器并且订阅初始化话题</button>
     </div>
  </div>
</div>

<div class="panel panel-default">
  <div class="panel-heading">动态订阅新话题</div>
  <div class="panel-body">
  	<div class="form-group">
       <div class="form-group">
    		<label for="exampleInputPassword1">话题</label>
    	<input type="text" class="form-control" id="anyTimetopics" placeholder="connectDesktop165241,connectDesktop165241">
    	<button onclick="subscribe()" type="botton" class="btn btn-primary">动态订阅</button>
  		</div>

     </div>
  </div>
</div>


<div class="panel panel-default">
  <div class="panel-heading">消息发送</div>
  <div class="panel-body">

    <form class="form-horizontal" onsubmit="return false">
      <div class="form-group">
        <label class="col-sm-2 control-label" for="inputEmail3">话题/用户</label>
        <div class="col-sm-10">
          <input type="text" id="totopic" class="form-control" value="ToUser://testuser">
          <p class="help-block">如果发送指定用户那么填写ToUser://+uid </p>
        </div>
      </div>

      <div class="form-group">
        <label class="col-sm-2 control-label" for="inputPassword3">SendType</label>
        <div class="col-sm-10">
          <input type="text" value="websocket" id="sendtype" class="form-control">
          <p class="help-block">目前仅支持websocket</p>
        </div>
      </div>

       <div class="form-group">
        <label class="col-sm-2 control-label" for="inputPassword3">内容</label>
        <div class="col-sm-10">
          <input type="text" value="" id="msg" class="form-control">
          <p class="help-block">发送给对方内容</p>
        </div>
      </div>

      <div class="form-group">
        <label class="col-sm-2 control-label" for="inputPassword3">附加数据</label>
        <div class="col-sm-10">
          <input type="text" value="" id="data" class="form-control">
          <p class="help-block">这里的必须是一个json，并且json的key和值都是string 如{"orderid":"1","price":"21"} 这个数据会原封不对转发出去</p>
        </div>
      </div>


      <div class="form-group">
        <div class="col-sm-offset-2 col-sm-10">
          <button class="btn btn-default" type="botton" onclick="postMsg()">发送</button>
        </div>
      </div>
    </form>




  </div>
</div>

</div>






</div>

<div class="col-md-4">
	<div class="panel panel-default">
  <div class="panel-heading">websocket输出</div>
  <div class="panel-body">
    	<div id="webconsolelog" style="height: 800px;overflow-y:auto;"></div>
  </div>
</div>
</div>


</div>

<script src="http://cdn.bootcss.com/jquery/1.12.4/jquery.js"></script>
<script src="https://cdn.bootcss.com/reconnecting-websocket/1.0.0/reconnecting-websocket.min.js"></script>
<script type="text/javascript">
	var socket
	function connectWebSocket(){
		url = $("#websocketUrl").val()
		url+="?uid="+$("#userid").val()
		url+="&topics="+$("#topics").val()

	 	socket = new ReconnectingWebSocket(url);
	    socket.onopen = function(event){
	    	appendtolog("成功连接"+url)
		}
		socket.onmessage = function(event){
				console.log(event)
				appendtolog("收到消息："+event.data)
			}

		socket.onclose = function(event){
				console.log(event)
				appendtolog("断开连接")
			}
	}

function appendtolog(str){
	var myDate = new Date();
	$("#webconsolelog").append("<p>"+myDate.toLocaleString()+"</p>")	
	$("#webconsolelog").append("<p>"+str+"</p>")
}

function postMsg(){
	data={}
	data.SendType = $("#sendtype").val()
	data.Msg = $("#msg").val()
	data.Topic = $("#totopic").val()
	data.data = $("#data").val()
	data.sendToUrl = $("#apiUrl").val()+"/send"

	 $.post("?ac=postMsg",data,function(result){

    	 appendtolog("发送消息2："+ JSON.stringify(data))
    	 appendtolog("发送结果："+result )

  	});
}

function subscribe(){
	data={}
	data.uid = $("#userid").val()
	data.topics = $("#anyTimetopics").val()
	data.sendToUrl = $("#apiUrl").val()+"/subscribe"

	 $.post("?ac=subscribe",data,function(result){
    	 appendtolog("订阅："+ $("#anyTimetopics").val())
    	 appendtolog("发送结果："+result )

  	});
}

</script>
 </body>
 </html>
 <?php }?>