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
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/metrumresearchgroup/fab/internal/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var collections []string

func configAdd(_ *cobra.Command, args []string) {
	// this will return like:
	// fab config init
	// osx:
	// /Users/<user>/Library/Application Support/fab/config.yml
	// linux:
	// /home/<user>/.config/fab/config.yml
	cfgPath, err := getDefaultConfigPath()
	if err != nil {
		log.Fatal(err)
		return
	}

	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	t := yaml.Node{}
	err = yaml.Unmarshal(cfgBytes, &t)
	newt, _ := config.AddPathsToCollections(t, collections, true, true)
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(&newt)
	fmt.Println(string(b.Bytes()))
	err = ioutil.WriteFile(cfgPath, b.Bytes(), 0644)
}

func newConfigAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add elements to a config",
		Run:   configAdd,
	}
	cmd.Flags().StringSliceVar(&collections, "collection", []string{}, "collection path to add")
	return cmd
}
