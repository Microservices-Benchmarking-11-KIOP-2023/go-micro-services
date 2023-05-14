# WSL
.PHONY: load-test

load-test:
	k6 run script.js --vus 10 --duration 5s --out json=report/test_5s_vus10.json --summary-export=report/summary_5s_vus10.json
	k6 run script.js --vus 10 --duration 10s --out json=report/test_10s_vus10.json --summary-export=report/summary_10s_vus10.json
	k6 run script.js --vus 10 --duration 100s --out json=report/test_100s_vus10.json --summary-export=report/summary_100s_vus10.json
	k6 run script.js --vus 10 --duration 1000s --out json=report/test_1000s_vus10.json --summary-export=report/summary_1000s_vus10.json

	k6 run script.js --vus 10 --duration 1000s --out json=report/test_vus10_duration1000s.json --summary-export=report/summary_vus10_duration1000s.json
	k6 run script.js --vus 20 --duration 1000s --out json=report/test_vus20_duration1000s.json --summary-export=report/summary_vus20_duration1000s.json
	k6 run script.js --vus 100 --duration 1000s --out json=report/test_vus100_duration1000s.json --summary-export=report/summary_vus100_duration1000s.json
	k6 run script.js --vus 1000 --duration 1000s --out json=report/test_vus1000_duration1000s.json --summary-export=report/summary_vus1000_duration1000s.json