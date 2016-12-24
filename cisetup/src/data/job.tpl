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
		<link rel="stylesheet" href="/css/bootstrap-select.min.css"/>
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

			.entry:not(:first-of-type)
			{
			margin-top: 10px;
			}

			.glyphicon
			{
			font-size: 12px;
			}

		</style>
	</head>
	<body>
		<script src="/js/jquery-3.1.1.min.js" type="text/javascript"></script>
		<script src="/js/bootstrap.min.js" type="text/javascript"></script>
		<script src="/js/bootstrap-select.min.js"></script>
		<script src="/js/bootstrap-formhelpers-phone.js"></script>
		<script>
			$(function()
			{
			$(document).on('click', '.btn-add', function(e)
			{
			e.preventDefault();

			var controlForm = $('.controls:first'),
			currentEntry = $(this).parents('.entry:first'),
			newEntry = $(currentEntry.clone()).appendTo(controlForm);

			newEntry.find('input').val('');
			controlForm.find('.entry:not(:last) .btn-add')
			.removeClass('btn-add').addClass('btn-remove')
			.removeClass('btn-success').addClass('btn-danger')
			.html('<span class="glyphicon glyphicon-minus"></span>');
			}).on('click', '.btn-remove', function(e)
			{
			$(this).parents('.entry:first').remove();

			e.preventDefault();
			return false;
			});
			});
		</script>

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
				<div class="panel-heading">Форма создания процесса</div>
				<div class="panel-body">
					<form role="form" autocomplete="off" action="/jobs/add" method="post">
						<div class="form-group{{if .JobName_err}} has-error{{end}}">
							<label for="JobName">Название</label>
							<input type="text" class="form-control input-sm" id="JobName" name="JobName" value="{{.JobName}}" />
							{{if .JobName_err}}<span class="help-block">{{.JobName_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="JobPrior">Приоритет</label>
							<select class="form-control selectpicker" id="JobPrior" name="JobPrior">
								<option value="1" {{if eq .JobPrior "1"}}selected{{end}}>Низкий</option>
								<option value="2" {{if eq .JobPrior "2"}}selected{{end}}>Стандартный</option>
								<option value="3" {{if eq .JobPrior "3"}}selected{{end}}>Высокий</option>
								<option value="9" {{if eq .JobPrior "9"}}selected{{end}}>Экстремальный</option>
							</select>
						</div>
						<div class="form-group{{if .JobCommit_err}} has-error{{end}}">
							<label for="JobCommit">Идентификатор изменения в системе контроля версий (имя ветки, коммит, или разница между коммитами). Для разницы коммитов два коммита должны быть указаны через запятую</label>
							<input type="text" class="form-control input-sm" id="JobCommit" name="JobCommit" value="{{.JobCommit}}"/>
							{{if .JobCommit_err}}<span class="help-block">{{.JobCommit_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="JobTask">Задача</label>
							<select class="form-control selectpicker" data-size="10" id="JobTask" name="JobTask"  data-live-search="true">
							    {{range .JobTask}}
							    <option value="{{index . 0}}" {{if eq (index . 1) "selected"}}selected{{end}}>{{index . 2}}</option>
							    {{end}}
							</select>
						</div>
						<div class="form-group">
							<label for="JobDescr">Дополнительное описание</label>
							<textarea class="form-control" rows="6" id="JobDescr" name="JobDescr">{{.JobDescr}}</textarea>
						</div>
						<button type="submit" class="btn btn-default">Отправить</button>
					</form>
				</div>
			</div>

		</div>
		<div class="panel-footer">BayZR Management Tool &copy; Alexey Berezhok</div>

	</body>
</html>