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
				<div class="panel-heading">Форма создания задачи</div>
				<div class="panel-body">
					<form role="form" autocomplete="off" action="/tasks/add" method="post">
						<div class="form-group{{if .TaskName_err}} has-error{{end}}">
							<label for="TaskName">Название</label>
							<input type="text" class="form-control input-sm" id="TaskName" name="TaskName" value="{{.TaskName}}" />
							{{if .TaskName_err}}<span class="help-block">{{.TaskName_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="TaskType">Тип результата</label>
							<select class="form-control selectpicker" id="TaskType" name="TaskType">
								<option value="1" {{if eq .TaskType "1"}}selected{{end}}>SonarQube</option>
								<option value="2" {{if eq .TaskType "2"}}selected{{end}}>Commit Check</option>
							</select>
						</div>
						<div class="form-group{{if .TaskGit_err}} has-error{{end}}">
							<label for="TaskGit">Команда клонирования</label>
							<input type="text" class="form-control input-sm" id="TaskGit" name="TaskGit" value="{{.TaskGit}}"/>
							{{if .TaskGit_err}}<span class="help-block">{{.TaskGit_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="TaskPackGs">Пакеты для сборки проекта</label>
							<div class="controls">
								<div class="entry input-group col-xs-3">
									<input class="form-control" id="TaskPackGs" name="TaskPackGs" type="text" />
									<span class="input-group-btn">
										<button class="btn btn-success btn-add" type="button">
											<span class="glyphicon glyphicon-plus"></span>
										</button>
									</span>
								</div>
							</div>
						</div>
						<div class="form-group">
							<label for="TaskPackGsEarl">Пакеты для сборки проекта ранее используемые</label>
							<select class="form-control selectpicker" data-size="10" multiple data-selected-text-format="count > 6" id="TaskPackGsEarl" name="TaskPackGsEarl"  data-live-search="true">
							    {{range .TaskPackGsEarl}}
							    <option value="{{index . 0}}" {{if eq (index . 1) "selected"}}selected{{end}}>{{index . 0}}</option>
							    {{end}}
							</select>
						</div>
						<div class="form-group">
							<label for="TaskCmds">Команды сборки (новая команда с новой строки)</label>
							<textarea class="form-control" rows="6" id="TaskCmds" name="TaskCmds">{{.TaskCmds}}</textarea>
						</div>
						<div class="form-group">
							<label for="TaskPerType">Тип периода</label>
							<select class="form-control selectpicker" id="TaskPerType" name="TaskPerType">
								<option value="0" {{if eq .TaskPerType "0"}}selected{{end}}>Ежеминутно</option>
								<option value="1" {{if eq .TaskPerType "1"}}selected{{end}}>Ежечасно</option>
								<option value="2" {{if eq .TaskPerType "2"}}selected{{end}}>Ежедневно</option>
								<option value="3" {{if eq .TaskPerType "3"}}selected{{end}}>Ежемесячно</option>
								<option value="4" {{if eq .TaskPerType "4"}}selected{{end}}>Один раз</option>
								<option value="5" {{if eq .TaskPerType "5"}}selected{{end}}>Без периода</option>
							</select>
						</div>
						<div class="form-group{{if .TaskPeriod_err}} has-error{{end}}">
							<label for="TaskPeriod">Время периода</label>
							<!-- Элемент HTML с id равным datetimepicker1 -->
							<input type="text" class="form-control bfh-phone" name="TaskPeriod" id="TaskPeriod" data-format="dd/dd/dddd dd:dd:dd" value="{{.TaskPeriod}}"/>
							{{if .TaskPeriod_err}}<span class="help-block">{{.TaskPeriod_err}}</span>{{end}}
						</div>
						<div class="form-group">
							<label for="TaskUsers">Кто может запускать</label>
							<select class="form-control selectpicker" data-size="10" multiple data-selected-text-format="count > 6" id="TaskUsers" name="TaskUsers"  data-live-search="true">
							    {{range .TaskUsers}}
							    <option value="{{index . 0}}" {{if eq (index . 1) "selected"}}selected{{end}}>{{index . 2}}</option>
							    {{end}}
							</select>
						</div>
						<div class="form-group{{if .TaskConfig_err}} has-error{{end}}">
							<label for="TaskConfig">Конфигурация для анализаторов кода</label>
							<textarea class="form-control" rows="6" id="TaskConfig" name="TaskConfig">{{.TaskConfig}}</textarea>
							{{if .TaskConfig_err}}<span class="help-block">{{.TaskConfig_err}}</span>{{end}}
						</div>
						<button type="submit" class="btn btn-default">Отправить</button>
					</form>
				</div>
			</div>

		</div>
		<div class="panel-footer">BayZR Management Tool &copy; Alexey Berezhok</div>

	</body>
</html>