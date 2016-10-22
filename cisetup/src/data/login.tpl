<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
   "http://www.w3.org/TR/html4/strict.dtd">
<html lang="en">
	<head>
		<meta charset="utf-8"/>
		<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
		<title>
			BayZR сервер задач для проверки кода
		</title>
		<link href="/css/bootstrap.min.css" rel="stylesheet" type="text/css"/>
		<style type="text/css">
			body {
			padding-top: 40px;
			padding-bottom: 40px;
			}

			.form-signin {
			max-width: 500px;
			padding: 15px;
			margin: 0 auto;
			}

			.form-signin {
			z-index: 2;
			}
		</style>
	</head>
	<body>
		<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
		<script src="/js/jquery-3.1.1.min.js" type="text/javascript">
		</script>
		<!-- Include all compiled plugins (below), or include individual files as needed -->
		<script src="/js/bootstrap.min.js" type="text/javascript">
		</script>

		<form role="form" class="form-signin" action="/" method="post">
			<h2 class="form-signin-heading">
				Добро пожаловать в управление заданиями анализатора кода BayZR
			</h2>
			{{if .IsError}}<div class="alert alert-danger">{{.ErrMSG}}</div>{{end}}
			<div class="form-group">
				<label for="InputLogin">логин</label>
				<input type="login" class="form-control" id="InputLogin" name="InputLogin" placeholder="Enter login" {{if .FormUser}}value="{{.FormUser}}"{{end}}/>
			</div>
			<div class="form-group">
				<label for="InputPasswd">Пароль</label>
				<input type="password" class="form-control" id="InputPasswd" name="InputPasswd" placeholder="Password"/>
			</div>
			<div class="form-group">
				<button type="submit" class="btn btn-primary" name="send" value="send">Отправить</button>
			</div>
			<div class="form-group">
				<button type="submit" class="btn btn-link btn-block" name="register" value="register">Зарегистрироваться</button>
			</div>
		</form>
	</body>
</html>