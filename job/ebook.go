package job

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ilovelili/dongfeng-jobs/core/model"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/micro/cli"
)

// ConvertEbookToPDF convert ebook to pdf
func ConvertEbookToPDF(ctx *cli.Context) int {
	operation := "ConvertEbookToPDF"
	ebooks, err := ebooksRepo.FindAll()
	if err != nil {
		errorLog(err.Error(), operation)
		return 1
	}

	for _, ebook := range ebooks {
		if err := convert(ebook); err != nil {
			errorLog(err.Error(), operation)
			return 1
		}
		if err := ebooksRepo.SetConverted(ebook); err != nil {
			errorLog(err.Error(), operation)
			return 1
		}
	}

	return 0
}

// MergeEbook merge ebook pdfs into one
func MergeEbook(ctx *cli.Context) int {
	operation := "MergeEbook"

	if err := merge(); err != nil {
		errorLog(err.Error(), operation)
		return 1
	}
	return 0
}

func convert(ebook *model.Ebook) (err error) {
	year, class, name, date := ebook.Pupil.Class.Year, ebook.Pupil.Class.Name, ebook.Pupil.Name, ebook.Date

	htmllocaldir := path.Join(config.Ebook.OriginDir, year, class, name, date)
	_, err = os.Stat(htmllocaldir)
	if err != nil && os.IsNotExist(err) {
		// original dir not exist, which is ok
		return nil
		// err = os.MkdirAll(htmllocaldir, os.ModePerm)
		// if err != nil {
		// 	return err
		// }
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	chromeDevTool := os.Getenv("CHROME_DEV_TOOL")
	if chromeDevTool == "" {
		chromeDevTool = "http://127.0.0.1:9222"
	}

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New(chromeDevTool)
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return
		}
	}
	defer devt.Close(ctx, pt)

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return
	}
	defer conn.Close() // Leaving connections open will leak memory.

	cli := cdp.NewClient(conn)
	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := cli.Page.DOMContentEventFired(ctx)
	if err != nil {
		return
	}
	defer domContent.Close()

	// Enable the runtime
	if err = cli.Runtime.Enable(ctx); err != nil {
		return
	}

	// Enable the network
	if err = cli.Network.Enable(ctx, network.NewEnableArgs()); err != nil {
		return
	}

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	if err = cli.Page.Enable(ctx); err != nil {
		return
	}

	// Create the Navigate arguments
	url := fmt.Sprintf("file://%s", path.Join(htmllocaldir, "index.html"))
	fmt.Println(url)
	navArgs := page.NewNavigateArgs(url)
	nav, err := cli.Page.Navigate(ctx, navArgs)
	if err != nil {
		return
	}

	// wait till image loaded
	time.Sleep(time.Duration(config.Ebook.ImageLoadTimeout) * time.Second)

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return
	}

	fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	// Print to PDF
	printToPDFArgs := page.NewPrintToPDFArgs().
		SetLandscape(false).
		SetPrintBackground(true).
		SetMarginTop(0).
		SetMarginBottom(0).
		SetMarginLeft(0).
		SetMarginRight(0).
		SetPaperWidth(config.Ebook.Width).
		SetPaperHeight(config.Ebook.Height)

	print, err := cli.Page.PrintToPDF(ctx, printToPDFArgs)
	if err != nil {
		return
	}

	pdfOutput := path.Join(htmllocaldir, "output.pdf")
	if err = ioutil.WriteFile(pdfOutput, print.Data, 0644); err != nil {
		return
	}

	fmt.Printf("Saved pdf: %s\n", pdfOutput)

	// move to dest dir
	pdfdestdir := path.Join(config.Ebook.PDFDestDir, year, class, name)
	_, err = os.Stat(pdfdestdir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(pdfdestdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.Rename(pdfOutput, path.Join(pdfdestdir, fmt.Sprintf("%s.pdf", ebook.Date)))
	return err
}

func merge() (err error) {
	// check if pdftk installed or not
	_, err = exec.LookPath("pdftk")
	if err != nil {
		return
	}

	filepathmap := make(map[string][]string)
	err = filepath.Walk(config.Ebook.MergeTargetDir, func(filepath string, info os.FileInfo, err error) error {
		// target
		if !info.IsDir() && path.Ext(info.Name()) == ".pdf" && info.Size() > 0 {
			key := path.Dir(filepath)
			// ignore the dest file
			if strings.Index(key, config.Ebook.MergeDestDir) > -1 {
				return nil
			}

			if paths, ok := filepathmap[key]; ok {
				filepathmap[key] = append(paths, filepath)
			} else {
				filepathmap[key] = []string{filepath}
			}
		}
		return nil
	})

	if err != nil {
		return
	}

	for dir := range filepathmap {
		segments := strings.Split(dir, "/")
		class, name := segments[len(segments)-2], segments[len(segments)-1]
		mergedestdir := path.Join(config.Ebook.MergeDestDir, class, name)
		// first clear the merge dir
		os.RemoveAll(mergedestdir)
	}

	for dir, filepaths := range filepathmap {
		// sort pdf by date
		sort.Strings(filepaths)
		// https://stackoverflow.com/questions/31467153/golang-failed-exec-command-that-works-in-terminal
		// cmdline := fmt.Sprintf("pdftk %s cat output merge.pdf", path.Join(filepath, "*.pdf"))
		pdffiles := strings.Join(filepaths, " ")
		cmdline := fmt.Sprintf("pdftk %s cat output %s", pdffiles, path.Join(dir, "merge.pdf"))
		args := strings.Split(cmdline, " ")
		cmd := exec.Command(args[0], args[1:]...)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}

		// move to dest
		segments := strings.Split(dir, "/")
		year, class, name := segments[len(segments)-3], segments[len(segments)-2], segments[len(segments)-1]
		mergedestdir := path.Join(config.Ebook.MergeDestDir, class, name)
		_, err = os.Stat(mergedestdir)
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(mergedestdir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		// class/name/year.pdf
		// 电子书_${this.currentName}_${this.currentClass}_${this.currentYear}学年.pdf
		err = os.Rename(path.Join(dir, "merge.pdf"), path.Join(mergedestdir, fmt.Sprintf("电子书_%s_%s_%s学年.pdf", name, class, year)))
		if err != nil {
			return
		}
	}

	// loop dest dir and merge again to generate the full year ebook
	destfilepathmap := make(map[string][]string)
	err = filepath.Walk(config.Ebook.MergeDestDir, func(filepath string, info os.FileInfo, err error) error {
		if !info.IsDir() && path.Ext(info.Name()) == ".pdf" {
			key := path.Dir(filepath)
			// ignore the target file
			// if strings.Index(key, config.Ebook.MergeTargetDir) > -1 {
			// 	return nil
			// }

			if paths, ok := destfilepathmap[key]; ok {
				destfilepathmap[key] = append(paths, filepath)
			} else {
				destfilepathmap[key] = []string{filepath}
			}
		}
		return nil
	})

	if err != nil {
		return
	}

	for dir, filepaths := range destfilepathmap {
		sort.Strings(filepaths)
		pdffiles := strings.Join(filepaths, " ")

		// move to dest
		segments := strings.Split(dir, "/")
		class, name := segments[len(segments)-2], segments[len(segments)-1]

		// 电子书_${this.currentName}_${this.currentClass}_全期间.pdf
		cmdline := fmt.Sprintf("pdftk %s cat output %s", pdffiles, path.Join(dir, fmt.Sprintf("电子书_%s_%s_全期间.pdf", name, class)))
		args := strings.Split(cmdline, " ")
		cmd := exec.Command(args[0], args[1:]...)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
	}

	return
}
