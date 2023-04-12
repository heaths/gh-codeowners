module github.com/{{param "github.owner"}}/{{param "name" (param "github.repo") "What is your project name?" | lowercase}}

go 1.18
