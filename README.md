# image-service

### Цель:

Написать сервис по хранению картинок на golang, который позволяет загружать и получать изображения по дате.

Требования:

- Реализовать HTTP API сервер: метод для добавления изображения, метод для получения изображений за определенный день.
- При помощи воркера раз в 10 минут конвертировать все jpg изображения в png.
- Контейнеризовать приложение с использованием Docker.
- Создать Makefile.


###  Примечание:
- Сохрание файлов происходит в уникальный UUID, можно из хедера брать имя файла и добавлять хэш, либо просто перезаписывать.
- Изображения скачиваются в зип файле, можно склеивать изображения, либо использовать stream writer.
- Во время запуска конвертора лочится вся папка с датами, если хотим оптимизировать, можем лочить, например, по часам, а не по дням.
- Можно добавить линтеры, хелсчеки, свагер, метрики, профилировщик...

### Проблемы:
- Сейчас скачивается только один img файл.

### Deploy:

```make run```

### Postman json collection (instead of generated swagger)
```
{
	"info": {
		"_postman_id": "26a33ece-5822-4eb7-87d0-5b38811d7ae0",
		"name": "Image-service",
		"schema": "https://schema.getpostman.com/json/collection/v2.0.0/collection.json",
		"_exporter_id": "24338277"
	},
	"item": [
		{
			"name": "http://localhost:8080/upload",
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "formdata",
					"formdata": [
						{
							"key": "image",
							"type": "file",
							"src": "TODO"
						}
					]
				},
				"url": "http://localhost:8080/upload"
			},
			"response": []
		},
		{
			"name": "http://localhost:8080/download",
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "http://localhost:8080/download?date=2024-09-15",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"download"
					],
					"query": [
						{
							"key": "date",
							"value": "2024-09-15"
						}
					]
				}
			},
			"response": []
		}
	]
}
```