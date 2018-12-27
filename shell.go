package main

import (
	"github.com/abiosoft/ishell"
	"github.com/kisasexypantera94/khalzam/musiclibrary"
	_ "github.com/lib/pq"
)

func main() {
	mLib, err := musiclibrary.Open()
	defer mLib.Close()
	if err != nil {
		panic(err)
	}

	shell := ishell.New()
	shell.Println("MusicLibrary interactive shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "index",
		Help: "index audiofile",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: index file ...")
			}

			for _, arg := range c.Args {
				err := mLib.Index(arg)
				if err != nil {
					c.Println(err)
					return
				}
				c.Println("Done")
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "indexdir",
		Help: "index directory",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: index dir ...")
			}

			for _, arg := range c.Args {
				err := mLib.IndexDir(arg)
				if err != nil {
					c.Println(err)
					return
				}
				c.Println("Done")
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: "delete audio from database",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: delete audio ...")
			}

			for _, arg := range c.Args {
				affected, err := mLib.Delete(arg)
				if err != nil {
					c.Println(err)
					return
				}
				if affected > 0 {
					c.Println("Done")
				} else {
					c.Println("Audio not found")
				}
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "recognize",
		Help: "recognize audiofile",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: recognize file ...")
			}

			for _, arg := range c.Args {
				res, err := mLib.Recognize(arg)
				if err != nil {
					c.Println(err)
					return
				}

				c.Println(res)
			}
		},
	})

	shell.Run()
}
