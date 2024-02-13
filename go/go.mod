module gonboat

go 1.22.0

require nvimboat v0.0.0-00010101000000-000000000000

require github.com/neovim/go-client v1.2.1

require (
	github.com/JohannesKaufmann/html-to-markdown v1.5.0 // indirect
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/net v0.19.0 // indirect
)

replace nvimboat => ./nvimboat
