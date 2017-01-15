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
			<div class="panel panel-success">
				<div class="panel-heading">Приветствие</div>
  				<div class="panel-body">
    				Добро пожаловать в систему запуска задач для статического анализа кода
  				</div>
			</div>
{{if or .ru_task .ru_result .ru_admin}}
			<div class="panel panel-success">
				<div class="panel-heading">Создание процесса</div>
  				<div class="panel-body">
    				<p>Процесс - это запущенная Задача с заданными параметрами. Результатом работы процесса 
    				является отчет о результатах проверки кода или сообщение об ошибке</p>
    				<p><b>Название</b> - название процесса, уникальное имя выполняемого задания</p>
    				<p><b>Приоритет</b> - приоритет выполнения процесса. При запросе анализа, процесс запускается не мгновенно,
    				а попадает в очередь процессов. Из каждой очечреди выбирается запрос на процесс с меньшим идентификатором. Порядок и очередность пересмотра очередей:
    				первой пересматривается очередь с приоритетом Экстремальный, потом Высокий, потом Стандартный и наконец - Низкий</p>
    				<p><b>Идентификатор изменения</b> - идентификатор изменения в системе контроля версий (имя ветки, коммит, или разница между коммитами). Для разницы коммитов два коммита должны быть указаны через запятую.
    				В качестве идентификатора может быть: имя ветки без remotes и origin(для задачити типа Ветка), или коммит или два коммита через запятую(для задачи типа Коммит)</p>
    				<p><b>Задача</b> - название задачи, параметры которой будут использованы</p>
    				<p><b>Дополнительное описание</b> - любой комментарий, может быть пустым</p>
  				</div>
			</div>
{{end}}
{{if or .ru_task .ru_admin}}
			<div class="panel panel-success">
				<div class="panel-heading">Создание задачи</div>
  				<div class="panel-body">
    				<p>Задача - это список параметров и действий будущего процесса</p>
    				<p><b>Название</b> - название задачи(name[:key[:version]]) key и version необязательны, при их отсутсвии они будут созданы из name</p>
    				<p><b>Тип результата</b> - CommitCheck не отправлять результат в SonarQube; SonarQube - отправить результат в SonarQube</p>
    				<p><b>Использовать коммит или ветку</b> - использовать в качестве идентификатора изменений коммит или имя ветки</p>
    				<p><b>Команда клонирования</b> - полная команда клонирования проекта. git clone ... </p>
    				<p><b>Пакеты для сборки проекта и Пакеты для сборки проекта ранее используемые</b> - список пакетов, которые должны быть доустановлены в окружение для успешного анализа или сборки проекта(если сборка необходима)</p>
    				<p><b>Команды сборки</b> - список команд, каждая с новой строки, которые необходимо выполнить чтоб получить корректные исходные коды для анализа.
    				Например для си-это может быть команда сборки проекта make VERBOSE=1 и т.д</p>
    				<p><b>Тип периода</b> - тип задачи, если Крон - то это периодическая задача, которую запускает сам сервис, иначе, Задача запускается вручную или по событию извне</p>
    				<p><b>Время периода</b> - описание крон задачи в крон формате <a href="https://godoc.org/github.com/robfig/cron">внешняя ссылка на описание формата</a></p>
    				<p><b>Кто может запускать</b> - список пользователей которым разрешен запуск, распространяется только на WEB API</p>
    				<p><b>Конфигурация для анализаторов кода</b> - bzr.conf со списокм анализаторов и настроек игнорирования. Больше информации по файлу: <a href="http://wiki.bayrepo.net/ru/bayzrdscr#bzrconf">внешняя ссылка</a></p>
    				<p><b>Файл результата</b> - имя файла результата проверки и путь к нему, если путь к нему не в корне git проекта, такое может быть если в пункте "Команды сборки" есть команды смены каталога</p>
    				<p><b>Ветка для периодических задач</b> - в какую ветку переключаться планировщику, при выполнеии задачи по крону (только для Cron задач)</p>
    				<p><b>Список проверяемых файлов</b> - в результирующий отчет могут попасть либо все файлы исходных кодов из проверяемого коммита(ветки) или только те файлы которые были найдены в разнице двух коммитов указанных при создании процесса через запятую</p>
    				<p><b>Команды выполняемые после аналитики</b> - команды для shell скрипта, которые выполняются после анализа, в shell скрипт передаются параметры - 1 - параметр число найденных ошибок, 2 - путь к файлу отчета, 3 - вывод команд.
    				Данные команды могут быть использованы для нотификации и пр.</p>
    				<p><b>Рабочий каталог в проекте</b> - каталог относительно корня git проекта, где будет запущен sonar-scanner</p>
    				<p><b>Команды выполняемые перед аналитикой от root</b> - это уже не установка пакетов, это может быть yum update или установка модулей питона pip install ... и т.д</p>
  				</div>
			</div>
{{end}}
		</div>	
		<div class="panel-footer">Утилита управления заданиями анализатора кода BayZR &copy; Alexey Berezhok</div>

	</body>
</html>