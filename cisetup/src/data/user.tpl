<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
   "http://www.w3.org/TR/html4/strict.dtd">
<html lang="en">
	<head>
		<meta name="generator" content=
				"HTML Tidy for Windows (vers 14 February 2006), see www.w3.org"/>
		<meta charset="utf-8"/>
		<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
		<title>
			BayZR сервер задач для проверки кода
		</title>
		<link href="/css/bootstrap.min.css" rel="stylesheet" type="text/css"/>
		<style>
			* {
			font-size: 12px;
			line-height: 1.428;
			}

			.center-panel {
			margin-top: 20px;
			margin-bottom: 20px;
			margin-left: 10%;
			margin-right: 10%;
			padding: 40px;
			}

		</style>
	</head>
	<body>
		<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
		<script src="/js/jquery-3.1.1.min.js" type="text/javascript"></script>
		<!-- Include all compiled plugins (below), or include individual files as needed -->
		<script src="/js/bootstrap.min.js" type="text/javascript"></script>

		<nav class="navbar navbar-default" role="navigation">
			<div class="container">
{{if or .ru_task .ru_result .ru_admin}}
				<form class="navbar-form navbar-left" role="search">
					<div class="form-group">
						<input type="text" class="form-control" placeholder="Search"/>
					</div>
					<button type="submit" class="btn btn-default">Find</button>
				</form>
{{end}}
				<ul class="nav navbar-nav">
{{if or .ru_task .ru_admin}}
					<li><a href="/tasks">Задания</a></li>
{{end}}
{{if or .ru_task .ru_result .ru_admin}}
					<li><a href="/procs">Процессы</a></li>
{{end}}
{{if or .ru_admin}}
					<li><a href="/users">Пользователи</a></li>
{{end}}
					<li><a href="/logout">Выход</a></li>
				</ul>
				<p class="navbar-text navbar-right">Вы вошли как <a href="/welcome">{{.User}}</a></p>
			</div>
		</nav>

		<div class="panel panel-default center-panel">
			<div class="panel panel-default">
				<div class="panel-heading">Форма параметров профиля</div>
				<div class="panel-body">
					<form role="form" action="/users/{{.UID}}" method="post">
						<div class="form-group">
							{{.User_m}}
						</div>
						<div class="form-group{{if .InputName_err}} has-error{{end}}">
							<label for="InputName">Имя</label>
							<input type="login" class="form-control input-sm" id="InputName" name="InputName" value="{{.Name}}" />
							{{if .InputName_err}}<span class="help-block">{{.InputName_err}}</span>{{end}}
						</div>
						<div class="form-group{{if .InputEmail1_err}} has-error{{end}}">
							<label for="InputEmail1">Email</label>
							<input type="email" class="form-control input-sm" id="InputEmail1" name="InputEmail1" value="{{.Email}}" />
							{{if .InputEmail1_err}}<span class="help-block">{{.InputEmail1_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="InputPassword1">Пароль</label>
							<input type="password" class="form-control input-sm" id="InputPassword1" name="InputPassword1" value=""/>
						</div>
						<div class="form-group{{if .InputPassword2_err}} has-error{{end}}">
							<label for="InputPassword2">Повтор пароля</label>
							<input type="password" class="form-control input-sm" id="InputPassword2" name="InputPassword2" value=""/>
							{{if .InputPassword2_err}}<span class="help-block">{{.InputPassword2_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="GroupInput">Группа</label>
							<select class="form-control" id="GroupInput" name="GroupInput">
							{{range .Groups}}
							<option value="{{index . 0}}" {{if eq (index . 0) $.Group_ID}}selected{{end}}>{{index . 1}}</option>
							{{end}}
							</select>
						</div>
						<button type="submit" class="btn btn-default">Отправить</button>
					</form>
				</div>
				<div>
				   <a href="/user/del/{{.UID}}">Удалить</a>
				</div>
			</div>
		</div>
		<div class="panel-footer">Утилита управления заданиями анализатора кода BayZR &copy; Alexey Berezhok</div>

	</body>
</html>