package cobra

import "strings"

// Command is a lightweight replacement for the real cobra.Command used in tests.
type Command struct {
    Use     string
    Aliases []string
    Short   string
    Long    string
    Version string
    Args    func(cmd *Command, args []string) error
    RunE    func(cmd *Command, args []string) error
    args    []string
    flags   *FlagSet
    children []*Command
}

func (c *Command) SetArgs(a []string) {
    c.args = a
    fs := c.Flags()
    for i := 0; i < len(a); i++ {
        arg := a[i]
        if strings.HasPrefix(arg, "--") {
            name := strings.TrimPrefix(arg, "--")
            val := "true"
            if i+1 < len(a) && !strings.HasPrefix(a[i+1], "--") {
                val = a[i+1]
                i++
            }
            fs.values[name] = val
            if p, ok := fs.strVars[name]; ok {
                *p = val
            }
            if b, ok := fs.boolVars[name]; ok {
                *b = val == "true"
            }
        }
    }
}

func (c *Command) Execute() error {
    if c.Args != nil {
        if err := c.Args(c, c.args); err != nil {
            return err
        }
    }
    if c.RunE != nil {
        return c.RunE(c, c.args)
    }
    return nil
}

func (c *Command) Flags() *FlagSet {
    if c.flags == nil {
        c.flags = &FlagSet{values: map[string]string{}}
    }
    return c.flags
}

// PersistentFlags mirrors Flags for this stub.
func (c *Command) PersistentFlags() *FlagSet { return c.Flags() }

// AddCommand adds subcommands to this command.
func (c *Command) AddCommand(cmds ...*Command) {
    c.children = append(c.children, cmds...)
}

// Commands returns the added subcommands.
func (c *Command) Commands() []*Command { return c.children }

// Name returns the command name which is its use line.
func (c *Command) Name() string {
    if idx := strings.Index(c.Use, " "); idx > 0 {
        return c.Use[:idx]
    }
    return c.Use
}

// FlagSet is a simple map based flag set supporting StringVarP, BoolVar and lookup.
type FlagSet struct{
    values map[string]string
    strVars map[string]*string
    boolVars map[string]*bool
}

func (f *FlagSet) StringVarP(p *string, name, shorthand, value, usage string) {
    if f.values == nil {
        f.values = map[string]string{}
    }
    if f.strVars == nil {
        f.strVars = map[string]*string{}
    }
    f.values[name] = value
    if p != nil {
        *p = value
        f.strVars[name] = p
    }
}

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
    if p != nil {
        *p = value
    }
    if f.values == nil {
        f.values = map[string]string{}
    }
    f.values[name] = ""
    if p != nil {
        if f.boolVars == nil { f.boolVars = map[string]*bool{} }
        f.boolVars[name] = p
    }
}

func (f *FlagSet) BoolP(name, shorthand string, value bool, usage string) *bool {
    b := value
    f.BoolVar(&b, name, value, usage)
    return &b
}

func (f *FlagSet) BoolVarP(p *bool, name, shorthand string, value bool, usage string) {
    f.BoolVar(p, name, value, usage)
}

func (f *FlagSet) Lookup(name string) *Flag {
    if f.values == nil {
        return nil
    }
    if _, ok := f.values[name]; ok {
        return &Flag{Name: name}
    }
    return nil
}

func (f *FlagSet) GetString(name string) (string, error) {
    if f.values == nil {
        return "", nil
    }
    return f.values[name], nil
}

type Flag struct{ Name string }

// MaximumNArgs returns a validator function for arguments.
func MaximumNArgs(n int) func(cmd *Command, args []string) error {
    return func(cmd *Command, args []string) error { return nil }
}
