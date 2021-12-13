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
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/metrumresearchgroup/fab/internal/config"
	"github.com/metrumresearchgroup/fab/internal/copier"
	"github.com/metrumresearchgroup/fab/internal/vcs"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gobuffalo/plush"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func generate(_ *cobra.Command, args []string) {
	gcfg, globalConfigPath, err := readConfig()
	if err != nil {
		log.Fatalf("Error reading config at %s: %v\n", globalConfigPath, err)
	}
	if len(gcfg.Collections)+len(gcfg.Templates) == 0 {
		log.Fatalf("Must define at least one collection or template in config %s\n", globalConfigPath)
	}
	err = copyTemplate(gcfg)
	if err != nil {
		log.Fatalf("Error copying template: %v", err)
	}
}

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a new project",
		Run:   generate,
	}

	return cmd
}

func copyTemplate(gcfg config.Config) error {
	// note as you read the internal implementation
	// from an error handling perspective there are two
	// schools of thought in conflict here - from one
	// side, you should just bubble up the errors to the
	// on the flip side, I want much more explicit logging
	// hence generally tried to use error messaging
	// but not fatally exit in this function
	rootDirs := gcfg.Templates

	for _, root := range gcfg.Collections {
		rfs := os.DirFS(root)
		setupFiles, err := fs.Glob(rfs, "**/_setup.yml")
		if err != nil {
			log.Warnf("error %v reading in template dirs at %v\n", err, root)
		}
		for _, d := range setupFiles {
			// the path to the setupFile will be relative to the root
			rootDirs = append(rootDirs,
				filepath.Join(root, filepath.Dir(d)),
			)
		}
	}

	if len(rootDirs) == 0 {
		return errors.New("no templates found")
	}
	ctx := plush.NewContext()
	// ctx.Set("configFile", true)
	// ctx.Set("cmd", "twinkle")
	// ctx.Set("module", "twinkle")
	// ctx.Set("organization", "metrumresearchgroup")
	// ctx.Set("repository", "twinkle")
	// ctx.Set("description", "shiny utilities")
	var qs = []*survey.Question{
		{
			Name: "root",
			Prompt: &survey.Select{
				Message: "Choose a project template:",
				Options: rootDirs,
			},
		},
		{
			Name:     "dest",
			Prompt:   &survey.Input{Message: "destination:"},
			Validate: survey.Required,
		},
		{
			Name:     "git",
			Validate: survey.Required,
			Prompt: &survey.Confirm{
				Message: "initialize as a git repo?",
				Default: true,
			},
		},
	}
	projSetup := struct {
		Root string
		Dest string
		Git  bool
	}{}
	err := survey.Ask(qs, &projSetup)
	if err != nil {
		log.Error("error asking questions")
		return err
	}
	projSetup.Dest, _ = homedir.Expand(projSetup.Dest)
	setupPath := filepath.Join(projSetup.Root, "_setup.yml")
	cfg, err := config.Read(setupPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Errorf("no setup file found at %s\n", setupPath)
			return err
		}
		log.Errorf("error reading setup file: %v\n", err)
		return err
	}
	qs = []*survey.Question{}
	for _, s := range cfg.Settings {
		if s.Type == "boolean" {
			qs = append(qs, &survey.Question{
				Name:     s.Name,
				Prompt:   &survey.Confirm{Message: s.Prompt},
				Validate: survey.Required,
			})
		} else {
			input := &survey.Input{
				Message: s.Prompt,
				Default: s.Default,
			}
			// overlay any global defaults - useful for known
			// common patterns like email or name where
			// could set in the global config
			for _, gs := range gcfg.Settings {
				if gs.Name == s.Name && gs.Default != "" {
					input.Default = gs.Default
				}
			}
			qs = append(qs, &survey.Question{
				Name:     s.Name,
				Prompt:   input,
				Validate: survey.Required,
			})
		}
	}
	// the questions to ask
	answers := make(map[string]interface{})
	// perform the questions
	err = survey.Ask(qs, &answers)
	if err != nil {
		log.Error("error asking questions")
		return err
	}
	for a, v := range answers {
		ctx.Set(a, v)
	}
	// fmt.Println(ctx.String())
	// one common error is that a template value doesn't exist
	// and plush will error with like:
	// "email": unknown identifier"
	err = copier.CopyDir(projSetup.Root, projSetup.Dest, ctx, answers)
	if err != nil {
		log.Errorf("error copying template %s: %v\n", projSetup.Root, err)
		log.Info("attempting to clean up...")
		err2 := os.RemoveAll(projSetup.Dest)
		if err2 != nil {
			log.Errorf("error removing directory %s: %v\n", projSetup.Dest, err2)
		}
		return err
	}
	log.Infof("new project created at: %s\n", projSetup.Dest)

	if projSetup.Git {
		err = vcs.UseGit(projSetup.Dest)
	}
	if err != nil {
		// we won't rollback/delete if not a git repo (for now)
		log.Errorf("error initializing git repo: %v\n", err)
	}
	return err
}
