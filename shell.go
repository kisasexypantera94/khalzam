package main

import (
	"github.com/abiosoft/ishell"
	"github.com/kisasexypantera94/khalzam/musiclibrary"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func main() {
	cfg := &musiclibrary.Config{
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASSWORD"),
		DBname:   os.Getenv("DBNAME"),
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
	}

	mLib, err := musiclibrary.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer mLib.Close()

	shell := ishell.New()
	shell.Println("MusicLibrary interactive shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "index",
		Help: "index audiofile",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: index file ...")
			}

			start := time.Now()
			for _, arg := range c.Args {
				err := mLib.Index(arg)
				if err != nil {
					c.Println(err)
					continue
				}
				c.Println("Done")
			}
			elapsed := time.Since(start)
			c.Printf("Finished in %s\n", elapsed)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "indexdir",
		Help: "index directory",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: index dir ...")
			}

			start := time.Now()
			for _, arg := range c.Args {
				err := mLib.IndexDir(arg)
				if err != nil {
					c.Println(err)
					continue
				}
				c.Println("Done")
			}
			elapsed := time.Since(start)
			c.Printf("Finished in %s\n", elapsed)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: "delete audio from database",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: delete audio ...")
			}

			start := time.Now()
			for _, arg := range c.Args {
				affected, err := mLib.Delete(arg)
				if err != nil {
					c.Println(err)
					continue
				}
				if affected > 0 {
					c.Println("Done")
				} else {
					c.Println("Audio not found")
				}
			}
			elapsed := time.Since(start)
			c.Printf("Finished in %s\n", elapsed)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "recognize",
		Help: "recognize audiofile",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: recognize file ...")
			}

			start := time.Now()
			for _, arg := range c.Args {
				res, err := mLib.Recognize(arg)
				if err != nil {
					c.Println(err)
					continue
				}

				c.Println(res)
			}
			elapsed := time.Since(start)
			c.Printf("Finished in %s\n", elapsed)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "recognizedir",
		Help: "recognize directory",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("usage: recognize dir ...")
			}

			start := time.Now()
			for _, arg := range c.Args {
				err := mLib.RecognizeDir(arg)
				if err != nil {
					c.Println(err)
					continue
				}
				c.Println("Done")
			}
			elapsed := time.Since(start)
			c.Printf("Finished in %s\n", elapsed)
		},
	})

	shell.Run()
}
