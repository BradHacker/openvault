package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := application.New(application.Options{
		Name:        "openvault",
		Description: "",
		Services: []application.Service{
			application.NewService(NewCoreService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Linux: application.LinuxOptions{
			ProgramName: "OpenVault",
		},
	})

	// app.RegisterService(application.NewService(NewCoreService()))

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "OpenVault",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHidden,
		},
		Frameless: true,
		// DisableResize:   false,
		BackgroundType:   application.BackgroundTypeTransparent,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		InitialPosition:  application.WindowCentered,
		Windows: application.WindowsWindow{
			BackdropType: application.Acrylic,
		},
		Linux: application.LinuxWindow{
			WindowIsTranslucent: true,
			// WebviewGpuPolicy: ,
			// Menu:                application.DefaultApplicationMenu(),
		},
		EnableDragAndDrop: true,
		URL:               "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// // Create application with options
	// err := wails.Run(&options.App{
	// 	Title:     "openvault",
	// 	Width:     1024,
	// 	Height:    768,
	// 	Frameless: true,
	// 	AssetServer: &assetserver.Options{
	// 		Assets: assets,
	// 	},
	// 	Windows: &windows.Options{
	// 		WindowIsTranslucent: true,
	// 	},
	// 	Mac: &mac.Options{
	// 		WindowIsTranslucent: true,
	// 	},
	// 	Linux: &linux.Options{
	// 		WindowIsTranslucent: true,
	// 	},
	// 	BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
	// 	OnStartup:        app.startup,
	// 	Bind: []interface{}{
	// 		app,
	// 	},
	// })

	// if err != nil {
	// 	println("Error:", err.Error())
	// }
}
