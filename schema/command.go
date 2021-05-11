package schema

type Command struct {
	builder *Builder

	CommandName       string
	CommandIndex      string
	CommandParameters []string
	CommandAlgorithm  string

	CommandReferences            []string
	CommandOnTable               string
	CommandOnDelete              string
	CommandOnUpdate              string
	CommandNotInitiallyImmediate bool
}

func NewCommand(builder *Builder) *Command {
	return &Command{
		builder:           builder,
		CommandParameters: make([]string, 0),
	}
}

func (c *Command) References(columns ...string) *Command {
	c.CommandReferences = columns
	return c
}

func (c *Command) On(table string) *Command {
	c.CommandOnTable = table
	return c
}

func (c *Command) OnDelete(action string) *Command {
	c.CommandOnDelete = action
	return c
}

func (c *Command) OnUpdate(action string) *Command {
	c.CommandOnUpdate = action
	return c
}

func (c *Command) NotInitiallyImmediate(value bool) *Command {
	c.CommandNotInitiallyImmediate = value
	return c
}

func (c *Command) Equal(name string) bool {
	return c.CommandName == name
}

func (c *Command) Name(name string) *Command {
	c.CommandName = name
	return c
}

func (c *Command) Index(name string) *Command {
	c.CommandIndex = name
	return c
}

func (c *Command) Columns(columns ...string) *Command {
	c.CommandParameters = columns
	return c
}

func (c *Command) Algorithm(algorithm string) *Command {
	c.CommandAlgorithm = algorithm
	return c
}

func (c *Command) Build() string {
	switch c.CommandName {
	case "index":
		return c.builder.compileKey(c, "index")
	case "unique":
		return c.builder.compileKey(c, "unique")
	case "primary":
		c.CommandIndex = ""
		return c.builder.compileKey(c, "primary key")
	case "spatialIndex":
		return c.builder.compileKey(c, "spatial index")
	case "dropTable":
		return c.builder.compileDropTable(c)
	case "dropColumn":
		return c.builder.compileDropColumn(c)
	case "dropIndex", "dropUnique", "dropSpatialIndex":
		return c.builder.compileDropIndex(c)
	case "dropPrimary":
		return c.builder.compileDropPrimary(c)
	case "dropTableIfExists":
		return c.builder.compileDropTableIfExists(c)
	case "rename":
		return c.builder.compileRename(c)
	case "create":
		return c.builder.compileCreateCommand(false)
	case "createIfNotExists":
		return c.builder.compileCreateCommand(true)
	case "add":
		return c.builder.compileAdd()
	case "change":
		return c.builder.compileChange()
	case "renameColumn":
		return c.builder.compileRenameColumn(c.CommandParameters[0], c.CommandParameters[1])
	case "foreign":
		return c.builder.compileForeign(c)
	case "dropForeign":
		return c.builder.compileDropForeign(c.CommandParameters[0])
	}

	return ""
}
