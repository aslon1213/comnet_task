package main
import _ "github.com/aslon1213/comnet_task/internal/app/initializers/tzinit"
import (
	
	"github.com/aslon1213/comnet_task/internal/pkg/app"
)

func main() {

	app := app.New()

	app.Run()

}
