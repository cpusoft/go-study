$(document).ready(function(){  
	//初始化查询来源
	loginQuery();
});

function loginQuery(){
	console.log("loginquery")
		$.ajax({
			url:'/loginQuery',
			type:'post',  
			dataType:'json',
			async : false ,
			cache : false ,
			data:"",
			success:function(data){  
				
					var username = data.username;
					var password = data.password;
					console.log(username)
					console.log(password)
					$("#username").attr("value",username);
					$("#password").attr("value",password);
				
			},  
			error:function(){  
			}  
		});  
}
	

function loginSubmit(){
	console.log("loginSubmit")
	var username = $("#username").val();
	var password = $("#password").val();
	console.log(username)
	console.log(password)
			
	$.ajax({
			url:'/loginSubmit',
			type:'post',  
			dataType:'json',
			async : false ,
			cache : false ,
			data:"username="+username+"&password="+password,
			success:function(data){  
				console.log(data)
					alert("login success")
					
				
			},  
			error:function(){  
			}  
	});  
}
