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
			BayZR job server
		</title>
		<link href="css/bootstrap.min.css" rel="stylesheet" type="text/css"/>
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
		<script src="js/jquery-3.1.1.min.js" type="text/javascript"></script>
		<!-- Include all compiled plugins (below), or include individual files as needed -->
		<script src="js/bootstrap.min.js" type="text/javascript"></script>

		<nav class="navbar navbar-default" role="navigation">
			<div class="container">
				<form class="navbar-form navbar-left" role="search">
					<div class="form-group">
						<input type="text" class="form-control" placeholder="Search"/>
					</div>
					<button type="submit" class="btn btn-default">Find</button>
				</form>
				<ul class="nav navbar-nav">
					<li><a href="/tasks">Задания</a></li>
					<li><a href="/procs">Процессы</a></li>
					<li><a href="/users">Пользователи</a></li>
					<li><a href="/logout">Выход</a></li>
				</ul>
				<p class="navbar-text navbar-right">Вы вошли как {{.User}}</p>
			</div>
		</nav>

		<div class="panel panel-default center-panel">
			<div class="panel panel-default">
				<div class="panel-heading">Форма параметров профиля</div>
				<div class="panel-body">
					<form role="form" action="/welcome" method="post">
						<div class="form-group">
							{{.User}}
						</div>
						<div class="form-group">
							<label for="InputName">Имя</label>
							<input type="email" class="form-control input-sm" id="InputName" name="InputName" value="{{.Name}}" />
						</div>
						<div class="form-group">
							<label for="InputEmail1">Email</label>
							<input type="email" class="form-control input-sm" id="InputEmail1" name="InputEmail1" value="{{.Email}}" />
						</div>
						<div class="form-group">
							<label for="InputPassword1">Пароль</label>
							<input type="password" class="form-control input-sm" id="InputPassword1" name="InputPassword1" value=""/>
						</div>
						<div class="form-group">
							<label for="InputPassword2">Повтор пароля</label>
							<input type="password" class="form-control input-sm" id="InputPassword2" name="InputPassword2" value=""/>
						</div>
						<div class="form-group">
							Группа: {{.Group}} Права:
							{{range .Rules}}
							<span class="label label-primary">{{.}}</span>
							{{end}}
						</div>
						<button type="submit" class="btn btn-default">Отправить</button>
					</form>
				</div>
			</div>

			<div class="row">
				<div class="col-md-4 text-center">
					<div class="panel panel-default">
						<div class="panel-body">
							Число проектов: {{.TaskCount}}
						</div>
					</div>
				</div>
				<div class="col-md-4 text-center">
					<div class="panel panel-default">
						<div class="panel-body">
							Число задач: {{.JobCount}}
						</div>
					</div>
				</div>
				<div class="col-md-4 text-center">
					<div class="panel panel-default">
						<div class="panel-body">
							Число пользователей: {{.UserCount}}
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="panel-footer">Утилита управления заданиями анализатора кода BayZR &copy; Alexey Berezhok</div>

	</body>
</html>