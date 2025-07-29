package color

// Green prints a formatted string. The real package adds color but this stub
// simply uses fmt.Printf without coloring.
import "fmt"

func Green(format string, a ...interface{}) {
    fmt.Printf(format+"\n", a...)
}
