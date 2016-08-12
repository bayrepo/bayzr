package visualmenu

import (
	"checker"
	"configparser"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var dialogs [][]string = [][]string{
	[]string{"Compillator", "compillator", ""},
	[]string{"File extension for check", "fileext", ""},
	[]string{"Custom definitions", "definit", ""},
	[]string{"Report file name", "repname", ""},
	[]string{"Replace compillator", "replace", ""},
	[]string{"List of checkers", "list", ""},
	[]string{"Output format", "fmt", ""},
	[]string{"Number of strings wrapped error", "wrap", ""},
}

var current_editable string = ""

type VisualMenu struct {
	conf *configparser.ConfigparserContainer
	plgs []*checker.PluginInfoDataContainer
}

func Make_VisualMenu(config *configparser.ConfigparserContainer, plg_lst []*checker.PluginInfoDataContainer) *VisualMenu {
	return &VisualMenu{config, plg_lst}
}

func (this *VisualMenu) Show() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(this.Layout)
	if err := this.keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func nextViewUp(g *gocui.Gui, v *gocui.View) error {
	err := cursorUp(g, v)
	if err != nil {
		return err
	}
	return nextView(g, v)
}

func nextViewDown(g *gocui.Gui, v *gocui.View) error {
	err := cursorDown(g, v)
	if err != nil {
		return err
	}
	return nextView(g, v)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	name := getLine(v)
	for _, val := range dialogs {
		if name == val[0] {
			_, err := g.SetViewOnTop(val[1])
			return err
		}
	}
	return nil
}

func getLine(v *gocui.View) string {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l
}

func getLineBegin(v *gocui.View) string {
	return strings.Join(strings.Split(strings.Trim(v.Buffer(), "\n\t "), "\n"), "")
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func jumpToDialog(g *gocui.Gui, v *gocui.View) error {
	line := getLine(v)
	for _, val := range dialogs {
		if val[0] == line {
			maxX, maxY := g.Size()
			if nv, err := g.SetView("msg", 35, maxY/2, maxX, maxY); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				fmt.Fprintln(nv, val[2])
				nv.Editable = true
				nv.Wrap = true
				nv.Title = val[0] + " (edit value, press Enter to save)"
				if err := g.SetCurrentView("msg"); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	new_value := getLineBegin(v)

	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if err := g.SetCurrentView("side"); err != nil {
		return err
	}

	if nv, err := g.View("side"); err != nil {
		return err
	} else {
		line := getLine(nv)
		for ind, val := range dialogs {
			if val[0] == line {
				dialogs[ind][2] = new_value
				if nvv, err2 := g.View(val[1]); err2 != nil {
					return err2
				} else {
					old_buffer := strings.Split(nvv.Buffer(), "\n")
					nvv.Clear()
					for _, val := range old_buffer {
						if strings.Contains(val, "-->") == false {
							fmt.Fprintln(nvv, val)
						} else {
							fmt.Fprintln(nvv, "--> "+new_value)
						}
					}
				}
			}
		}
	}

	return nil
}

func bzrd_exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}

func CopyDir(source string, dest string) (err error) {

	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return &CustomError{"Source is not a directory"}
	}

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return &CustomError{"Destination already exists"}
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {

			err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return
}

type CustomError struct {
	What string
}

func (e *CustomError) Error() string {
	return e.What
}

func (this *VisualMenu) saveConfig(g *gocui.Gui, v *gocui.View) error {
	var config configparser.ConfigparserContainer = *this.conf

	value, err := strconv.ParseInt(strings.Trim(dialogs[7][2], " \n\t"), 10, 64)
	if err != nil {
		return err
	}
	config.SetFields(
		configparser.SplitOwn(strings.Trim(dialogs[0][2], " \n\t")),
		configparser.SplitOwn(strings.Trim(dialogs[1][2], " \n\t")),
		configparser.SplitOwn(strings.Trim(dialogs[2][2], " \n\t")),
		strings.Trim(dialogs[3][2], " \n\t"),
		(strings.ToLower(strings.Trim(dialogs[4][2], " \n\t")) == "on"),
		configparser.SplitOwn(strings.Trim(dialogs[5][2], " \n\t")),
		strings.ToLower(strings.Trim(dialogs[6][2], " \n\t")),
		value)

	config.WriteToFile("bzr.conf")

	if exst, _ := bzrd_exists("bzr.d"); exst == false {
		if exst_etc, err_dir := bzrd_exists("/etc/bzr.d"); exst_etc == true && err_dir == nil {
			if currentPath, err := os.Getwd(); err == nil {
				err = CopyDir("/etc/bzr.d", currentPath+"/bzr.d")
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return gocui.ErrQuit
}

func (this *VisualMenu) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, jumpToDialog); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, nextViewDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, nextViewUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyCtrlS, gocui.ModNone, this.saveConfig); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	return nil
}

func showDialog(g *gocui.Gui, textv string, value string, name string) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(name, 35, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, textv)
		fmt.Fprintln(v, "Current value for project:")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "--> "+value)
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "Press ENTER to edit")
		fmt.Fprintln(v, "Press CTRL+C to exit")
		fmt.Fprintln(v, "Press CTRL+S to save values")
		v.Editable = false
		v.Wrap = true
		for ind, val := range dialogs {
			if val[1] == name {
				dialogs[ind][2] = value
			}
		}
		return err
	}
	return nil
}

func (this *VisualMenu) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 35, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		for _, val := range dialogs {
			fmt.Fprintln(v, val[0])
		}
	}

	if err := showDialog(g,
		"Option that allows to set list of extensions which will be parsed by bayzr.",
		strings.Join(this.conf.GetFilesList(), ", "),
		dialogs[1][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if err := showDialog(g,
		"Option that allows to set list of custom definitions and includes which will be added to command line for analytic tool.",
		strings.Join(this.conf.GetGlobalDefs(), ", "),
		dialogs[2][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	output, _, wrapstrings, outputfile := this.conf.GetReport()
	if err := showDialog(g,
		"Option that allows to set report file name for report.",
		outputfile,
		dialogs[3][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	replacer := this.conf.Replacer()
	replayser_str := "No"
	if replacer {
		replayser_str = "Yes"
	}
	if err := showDialog(g,
		"Option that allows to set - will it be replace gcc or c++ compillator in make for more precicional checking.",
		replayser_str,
		dialogs[4][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	available_plugins := []string{}
	for _, plg_item := range this.plgs {
		available_plugins = append(available_plugins, plg_item.GetNameId()+": "+plg_item.GetName())
	}
	if err := showDialog(g,
		"Option that allows to set list of checkers.\n"+strings.Join(available_plugins, "\n"),
		strings.Join(this.conf.GetListOfPlugins(), ", "),
		dialogs[5][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if err := showDialog(g,
		"Option that allows to set output format of result.\nAvailable:\n custom, txt, html\n",
		output,
		dialogs[6][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if err := showDialog(g,
		"Option that allows to set number of wrap strings for error string wrapping for report.",
		strconv.FormatInt(wrapstrings, 10),
		dialogs[7][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	if err := showDialog(g,
		"Option that allows to set list of compillators which will be parsed by bayzr.",
		strings.Join(this.conf.GetCompillatorsList(), ", "),
		dialogs[0][1]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.SetCurrentView("side"); err != nil {
			return err
		}
	}

	return nil
}
