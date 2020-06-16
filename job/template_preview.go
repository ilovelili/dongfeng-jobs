package job

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/ilovelili/dongfeng-jobs/core/model"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/micro/cli"
)

// ConvertTemplatePreviewToPDF convert template preview to pdf
func ConvertTemplatePreviewToPDF(ctx *cli.Context) int {
	operation := "ConvertTemplatePreviewToPDF"
	templatePreviews, err := ebooksRepo.FindAllTemplatePreviews()
	if err != nil {
		errorLog(err.Error(), operation)
		return 1
	}

	for _, templatePreview := range templatePreviews {
		if err := convertTemplatePreview(templatePreview); err != nil {
			errorLog(err.Error(), operation)
			return 1
		}
		if err := ebooksRepo.SetTemplatePreviewsConverted(templatePreview); err != nil {
			errorLog(err.Error(), operation)
			return 1
		}
	}

	return 0
}

func convertTemplatePreview(templatePreview *model.TemplatePreview) (err error) {
	htmllocaldir := path.Join(config.Ebook.OriginDir, "templatePreview", templatePreview.Name)
	_, err = os.Stat(htmllocaldir)
	if err != nil && os.IsNotExist(err) {
		// original dir not exist, which is ok
		return nil
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
	pdfdestdir := path.Join(config.Ebook.PDFDestDir, "templatePreview")
	_, err = os.Stat(pdfdestdir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(pdfdestdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.Rename(pdfOutput, path.Join(pdfdestdir, fmt.Sprintf("%s.pdf", templatePreview.Name)))
	return err
}
