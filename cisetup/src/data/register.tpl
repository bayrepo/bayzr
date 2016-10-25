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
			Регистрация для BayZR сервера задач
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

		<div class="panel panel-default center-panel">
			<div class="panel panel-default">
				<div class="panel-heading">Форма параметров профиля</div>
				<div class="panel-body">
					<form role="form" action="/register" method="post">
						<div class="form-group{{if .InputLogin_err}} has-error{{end}}">
							<label for="InputLogin">Логин</label>
							<input type="login" class="form-control input-sm" id="InputLogin" name="InputLogin" value="{{.InputLogin}}" />
							{{if .InputLogin_err}}<span class="help-block">{{.InputLogin_err}}</span>{{end}}
						</div>
						<div class="form-group{{if .InputName_err}} has-error{{end}}">
							<label for="InputName">Имя</label>
							<input type="login" class="form-control input-sm" id="InputName" name="InputName" value="{{.InputName}}" />
							{{if .InputName_err}}<span class="help-block">{{.InputName_err}}</span>{{end}}
						</div>
						<div class="form-group{{if .InputEmail1_err}} has-error{{end}}">
							<label for="InputEmail1">Email</label>
							<input type="email" class="form-control input-sm" id="InputEmail1" name="InputEmail1" value="{{.InputEmail1}}" />
							{{if .InputEmail1_err}}<span class="help-block">{{.InputEmail1_err}}</span>{{end}}
						</div>
						<div class="form-group{{if .InputPassword1_err}} has-error{{end}}">
							<label for="InputPassword1">Пароль</label>
							<input type="password" class="form-control input-sm" id="InputPassword1" name="InputPassword1" value="{{.InputPassword1}}"/>
							{{if .InputPassword1_err}}<span class="help-block">{{.InputPassword1_err}}</span>{{end}}
						</div>
						<div class="form-group{{if .InputPassword2_err}} has-error{{end}}">
							<label for="InputPassword2">Повтор пароля</label>
							<input type="password" class="form-control input-sm" id="InputPassword2" name="InputPassword2" value="{{.InputPassword2}}"/>
							{{if .InputPassword2_err}}<span class="help-block">{{.InputPassword2_err}}</span>{{end}}
						</div>
						<button type="submit" class="btn btn-default">Отправить</button>
					</form>
				</div>
			</div>
		</div>
		<div class="panel-footer">Утилита управления заданиями анализатора кода BayZR &copy; Alexey Berezhok</div>

	</body>
</html>