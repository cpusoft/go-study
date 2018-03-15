<html>
<head>
<title></title>
</head>
<body>
<form action="/login" method="post">
	用户名:<input type="text" name="username" value="{{.username}}">
	密码:<input type="password" name="password" value="{{.password}}">
	<input type="submit" value="登录">
</form>
</body>
</html>