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

                <a href="/tasks/add"><span class="glyphicon glyphicon-trash"></span>Создать новую задачу</a>
				<table class="table table-striped">
					<tr>
						<th>#</th>
						<th>Название</th>
						<th>Тип</th>
						<th>Хранилище</th>
						<th>Пакеты</th>
						<th>Команды</th>
						<th>Период</th>
						<th>Доступ</th>
						<th>Конфигурация</th>
						<th>Кто создал</th>
						<th>Действие</th>
					</tr>
{{range .Tasks}}
					<tr>
						<td>{{index . 0}}</td>
						<td>{{index . 1}}</td>
						<td>
						{{if eq (index . 2) "1"}}SonarQube{{end}}
						{{if eq (index . 2) "2"}}Commit Check{{end}}
						</td>
						<td>{{index . 3}}</td>
						<td>{{index . 4}}</td>
						<td>{{index . 5}}</td>
						<td>
						{{if eq (index . 6) "0"}}Ежеминутно{{end}}
						{{if eq (index . 6) "1"}}Ежечасно{{end}}
						{{if eq (index . 6) "2"}}Ежедневно{{end}}
						{{if eq (index . 6) "3"}}Ежемесячно{{end}}
						{{if eq (index . 6) "4"}}Один раз{{end}}
						{{if eq (index . 6) "5"}}Без периода{{end}}
						</td>
						<td>{{index . 10}}</td>
						<td>{{index . 9}}</td>
						<td>{{index . 8}}</td>
						<td>
							<a href="/task/{{index . 0}}">Изменить</a>
							<a href="/taskdel/{{index . 0}}">Удалить</a>
						</td>
					</tr>
{{end}}
				</table>
		</div>	
		<div class="panel-footer">Утилита управления заданиями анализатора кода BayZR &copy; Alexey Berezhok</div>

	</body>
</html>