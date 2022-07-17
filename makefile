gen-report:
	@cat example/crontab_service | go run ./cronv/main.go -d 7d -t "Cron Timeline (Prod)" --from-time="00:00" --from-date="2022/01/01" && open crontab.html

open:
	@cat example/crontab_service | go run ./cronv/main.go -d 7d -t "Cron Timeline (Prod)" --from-time="00:00" --from-date="2022/01/01"