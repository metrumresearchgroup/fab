/*
Copyright © 2021 Metrum Research Group <developers@metrumrg.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type settings struct {
	// strict mode will prevent the following:
	// - will check for path existence
	strict bool
	// logrus log level
	loglevel string
}

var cfg settings

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "`fab`ricate new projects in a `fab`ulous way",
}

func Execute(version string, commit string, date string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{printf \"%s\\n\" .Version}}")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newDebugCmd())
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newConfigCmd())

	// using viper so can take advantage of the casting and lookup capabilities of viper
	// even if don't need some of the more advanced functionality
	rootCmd.PersistentFlags().Bool("no-strict", false, "no strict mode")
	viper.BindPFlag("no-strict", rootCmd.PersistentFlags().Lookup("no-strict"))
	rootCmd.PersistentFlags().String("loglevel", "info", "log level")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	cobra.OnInitialize(initConfig)

}

func initConfig() {
	cfg.strict = !viper.GetBool("no-strict")
	cfg.loglevel = viper.GetString("loglevel")
	setLogLevel(cfg.loglevel)
}
