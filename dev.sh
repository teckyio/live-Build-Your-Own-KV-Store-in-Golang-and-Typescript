set -e
echo "set alice apple
set bob banana
get alice
del alice
get alice
rename bob alice
get alice
exit" |
	go run main.go
