package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/UnwrittenFun/ns/config"
	"github.com/UnwrittenFun/ns/util"
	"github.com/dixonwille/wmenu"
	homedir "github.com/mitchellh/go-homedir"
	wlog "gopkg.in/dixonwille/wlog.v2"
)

var cui *wlog.ColorUI
var cfg *config.Config

func main() {
	var ui wlog.UI
	ui = wlog.New(os.Stdin, os.Stdout, os.Stdout)
	ui = wlog.AddConcurrent(ui)
	cui = wlog.AddColor(
		wlog.Green,
		wlog.BrightRed,
		wlog.None,
		wlog.None,
		wlog.Cyan,
		wlog.None,
		wlog.None,
		wlog.BrightBlue,
		wlog.Yellow,
		ui,
	)

	homedir, err := homedir.Dir()
	if err != nil {
		cui.Error("Failed to locate home directory")
		cui.Error(err.Error())
		return
	}

	nsHome := filepath.Join(homedir, ".ns")
	if err := os.MkdirAll(nsHome, os.ModePerm); err != nil {
		cui.Error("Failed to locate create .ns home directory")
		cui.Error(err.Error())
		return
	}

	cfg = config.New(nsHome, "ns.json")

	files, err := ioutil.ReadDir(nsHome)
	if err != nil {
		cui.Error("Failed to walk `.ns` directory")
		cui.Error(err.Error())
		return
	}

	var npmrcs []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		ext := filepath.Ext(fileName)
		if ext != ".npmrc" {
			continue
		}
		name := strings.TrimSuffix(fileName, ext)
		npmrcs = append(npmrcs, name)
	}

	setupMenuAndRun(cfg.GetString("user", npmrcs[0]), npmrcs, func(opt []wmenu.Opt) error {
		npmrcPath := filepath.Join(nsHome, opt[0].Text+".npmrc")
		if err := util.CopyFile(filepath.Join(homedir, ".npmrc"), npmrcPath, os.ModePerm); err != nil {
			return err
		}

		cfg.Set("user", opt[0].Text)
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Println()
		cui.Success("Switched to user " + opt[0].Text)
		return nil
	})
}

func setupMenuAndRun(current string, options []string, handler func(opt []wmenu.Opt) error) {
	menu := wmenu.NewMenu("Switch account to:")
	menu.ClearOnMenuRun()
	menu.LoopOnInvalid()
	menu.AddColor(cui.OutputFGColor, cui.AskFGColor, wlog.None, cui.ErrorFGColor)
	menu.SetDefaultIcon("* ")

	for _, option := range options {
		menu.Option(option, nil, current == option, nil)
	}
	menu.Action(handler)

	menu.Option("~~ Exit ~~", nil, false, func(opt wmenu.Opt) error {
		os.Exit(0)
		return nil
	})

	err := menu.Run()
	if err != nil {
		cui.Error(err.Error())
		return
	}
}
