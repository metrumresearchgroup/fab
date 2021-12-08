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
	"fab/internal/config"

	"github.com/adrg/xdg"
)

func getDefaultConfigPath() (string, error) {
	return xdg.ConfigFile("fab/config.yml")
}

func readConfig() (config.Config, string, error) {
	// for now don't check on error until consider what type of
	// logging - this should be completely optional anyway
	globalConfigPath, err := xdg.SearchConfigFile("fab/config.yml")
	if err != nil {
		// this message is and print out all the paths that were searched
		// for example, by mangling the name to `onfig.yml` to see what the
		// error was, got the following result:
		//
		// could not locate `onfig.yml` in any of the following paths: /Users/devinp/Library/Application Support/fab, /Users/devinp/Library/Preferences/fab, /Library/Application Support/fab, /Library/Preferences/fab
		// exit status 1
		return config.Config{}, "", err
	}
	cfg, err := config.Read(globalConfigPath)
	return cfg, globalConfigPath, err
}
