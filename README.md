## Access Logs Stats

[![Go Report Card](https://goreportcard.com/badge/blackbass1988/access_logs_stats)](https://goreportcard.com/report/github.com/blackbass1988/access_logs_stats)

The purpose of the application is to make a piece that could process "very many" access logs and build statistics on them

GOLang setup

------------

1) install via yum/apt/brew install go...
2) install via [gvm](https://github.com/moovweb/gvm#installing) (аналог rvm)
3) install via [docker](https://www.docker.com/products/overview)
4) install from [official site](https://golang.org/dl/)


install to GOPATH
---------------------

if you don't have $GOPATH setup it
```
export GOPATH=~/go
mkdir $GOPATH/{src,bin,pkg}
```

```
git clone https://github.com/blackbass1988/access_logs_stats $GOPATH/src/github.com/blackbass1988/access_logs_stats
```

Build
------

common case
```
make
```

UPX
--------

после того, как добро собралось, можно пакнуть с использованием https://upx.github.io
не стал добалять в однострочник, потому что сначала надо поставить upx.
Пусть это остается на совести

В результате использования флагов и упаковщика получается уменьшить размер бинарника с 3.6MB до 845KB

```
upx access_logs_stats
```

usage
-------------

```
./access_logs_stats -c config.json
```

```
./access_logs_stats -c config.yaml
```

[config.json example](config.json.example)

[config.yaml example](config.yaml.example)

configuration
---------


|field|description|
|----|------|
|*input*| это точка, откуда будут читаться. Здесь может быть как файл,так и пайп, например. *Experimental: syslog:udp::515/nginx*|
|*regexp*|глобальное регулярное выражение, которое нужно, чтобы выделить _поля_ для последующих вычислений|
|*period*|период, раз во сколько отправлять статистику в output. Валидные значения единиц измерения - "ns", "us" (или "µs"), "ms", "s", "m", "h".|
|*filters*|перечисление фильтров, по которым будут считаться метрики. Таким образом можно в отдельности считать метрики по каждому фильтру. Описание формата фильтра описано ниже|
|*output*|перечисление методов отправки результатов. У каждого отправителя  может быть своя настройка. Список доступных отправителей и способе их настройки описан ниже|
|*template_vars*|объект переменных, которые можно поместить в output.template или input в формате ${variableName}|

*input*

one of:

* file
* syslog
* stdin:nowait

**Filter**

|field|description|
|----|------|
|*filter*| регулярное выражение, описывающее, какие строки должны попасть под фильтр |
|*prefix*| префикс, который будет у ключа в output. |
|*items*| массив. перечисление метрик, которые надо посчитать и отправить в output |
|*items[].field*| названия поля. Соответствует полям из глобального регулярного выражения _regexp_ |
|*metrics*| перечисление метрик, которые надо посчитать для поля _field_|

**Output**

На данный момент доступно 2 отправщика: console и zabbix
console не имеет настроек, заббикс имеет следующие настройки
zabbix_host - хост сервера zabbix, 
zabbix_port - порт сервера zabbix, 
host - имя хоста, которым будет представляться приложение при отправке результатов

общий формат отправщика:

```
{
"type": "output_name",
"settings": {"output_config1":"output_config_value1", template: "${metric}.${field}"}
}
```

в случае отправщика console надо оставить объект settings пустым (settings:{})

**Формат ключа в отправщик**
по умолчанию формат следующий:

${field}.${metric}

если указан prefix у фильтра, то он ВСЕГДА добавляется в ${field}

можно поменять формат вывода, переопределив параметр template в свойстве settings конкретного output

у двух разных output может быть два разных template

**Список доступных операций со счетчиками (counts):**

* cps_{val} - кол-во элементов по уникальному значению _{val}_ в секунду для поля _field_
* uniq - кол-во уникальных значений за период
* percentage_{val} - процент по уникальному _{val}_ за съем для поля _field_

**Список доступных групповых операци (aggregated):**

Сохраяняет все значения из поля (с плавающей запятой)
 и позволяет применить следующие операции:

* min - минимальное значение по полю
* max - максимальное значение по полю
* avg - среднее значение по полю
* ips (items per second)
* len (кол-во элементов в группе), 
* cent_{N} - посчитать N-ый перцентиль



Run
------

```
./access_logs_stats -c config.json
```


todo
-----------------

make english doc

make tests for sender 

make normal syslog parser and remove regular expressions

make conf.d/*.json for multiple instances of app

make normal exit after one tick

re/libpcre.go getNamedGroupsFromExpression make parser 
